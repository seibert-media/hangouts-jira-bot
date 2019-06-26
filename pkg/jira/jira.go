package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"cloud.google.com/go/pubsub"
	goJira "github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-hangouts"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

// JIRA connection
type JIRA struct {
	*goJira.Client
	hangouts *hangouts.Hangouts
	baseURL  string
}

// New JIRA connection
func New(ctx context.Context, url, username, password string, hangouts *hangouts.Hangouts) *JIRA {
	ctx = log.WithFields(ctx,
		zap.String("component", "jira"),
		zap.String("url", url),
		zap.String("username", username),
	)

	tp := goJira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := goJira.NewClient(tp.Client(), url)
	if err != nil {
		log.From(ctx).Error("connecting jira", zap.Error(err))
	}
	log.From(ctx).Info("connected jira")

	return &JIRA{
		Client:   client,
		hangouts: hangouts,
		baseURL:  url,
	}
}

// Callback for PubSub messages
func (j *JIRA) Callback(ctx context.Context, m *pubsub.Message) {
	var msg hangouts.Event
	err := json.Unmarshal(m.Data, &msg)
	if err != nil {
		log.From(ctx).Error("parsing event", zap.Error(err))
		j.hangouts.Send(ctx, msg.Space.Name, j.BuildError(errors.Wrap(err, "parsing event"), msg.Message.Thread.Name))
		return
	}

	log.From(ctx).Debug(
		"receiving event",
		zap.String("type", msg.Type),
		zap.String("space", msg.Space.Name),
		zap.String("user", msg.User.Email),
		zap.String("msg", msg.Message.Text),
	)

	m.Ack()

	regex := regexp.MustCompile("[A-Za-z]+-[1-9][0-9]*")
	issues := regex.FindAllString(msg.Message.Text, -1)
	for _, issue := range issues {
		log.From(ctx).Debug("handling issue", zap.String("id", issue))

		m, err := j.BuildMessage(ctx, issue, msg)
		if err != nil {
			log.From(ctx).Warn("building message", zap.Error(err))
			mb, _ := json.Marshal(m)
			log.From(ctx).Debug("message", zap.String("content", string(mb)))
			j.hangouts.Send(ctx, msg.Space.Name, j.BuildError(errors.Wrap(err, "build card failed"), msg.Message.Thread.Name))
			continue
		}

		err = j.hangouts.Send(ctx, msg.Space.Name, m)
		if err != nil {
			log.From(ctx).Error("sending message", zap.Error(err))
			mb, _ := json.Marshal(m)
			log.From(ctx).Debug("message", zap.String("content", string(mb)))
			j.hangouts.Send(ctx, msg.Space.Name, j.BuildError(errors.Wrap(err, "send card error"), msg.Message.Thread.Name))
			continue
		}
	}
}

// BuildMessage for issue
func (j *JIRA) BuildMessage(ctx context.Context, issue string, msg hangouts.Event) (*hangouts.Message, error) {
	i, _, err := j.Client.Issue.Get(issue, nil)
	if err != nil {
		log.From(ctx).Debug("jira issue fetch error", zap.String("issue", issue), zap.Error(err))
		return nil, err
	}

	assignee := "Unassigned"
	if i.Fields.Assignee != nil {
		assignee = i.Fields.Assignee.Name
	}

	m := hangouts.NewMessage().InThread(msg.Message.Thread.Name).
		WithCard(
			hangouts.NewCard().WithHeader(
				fmt.Sprintf("%s: %s", issue, i.Fields.Summary),
				"",
				"https://storage.googleapis.com/bot-icons/jira-app-icon.png",
				"AVATAR",
			).WithSection(
				hangouts.NewSection("").WithWidget(hangouts.NewWidget().
					WithKeyValue(hangouts.NewKeyValue("Status", i.Fields.Status.Name, false)),
				).WithWidget(hangouts.NewWidget().
					WithKeyValue(hangouts.NewKeyValue("Assignee", assignee, false)),
				),
			).WithSection(
				hangouts.NewSection("").WithWidget(hangouts.NewWidget().
					WithButton(hangouts.NewTextLinkButton(
						"Open Ticket",
						fmt.Sprintf("%s/browse/%s", j.baseURL, issue)),
					),
				),
			),
		)

	return m, nil
}

// BuildError message
func (j *JIRA) BuildError(e error, thread string) *hangouts.Message {
	return hangouts.NewMessage().InThread(thread).WithText(fmt.Sprintf("Error: %s", e.Error()))
}
