package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"cloud.google.com/go/pubsub"
	goJira "github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/playnet-public/libs/log"
	"github.com/seibert-media/go-hangouts"
	"go.uber.org/zap"
)

// JIRA connection
type JIRA struct {
	*log.Logger
	*goJira.Client
	hangouts *hangouts.Hangouts
	baseURL  string
}

// New JIRA connection
func New(log *log.Logger, url, username, password string, hangouts *hangouts.Hangouts) *JIRA {
	log = log.WithFields(
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
		log.Error("jira connection error", zap.Error(err))
	}
	log.Info("jira connected")

	return &JIRA{
		Logger:   log,
		Client:   client,
		hangouts: hangouts,
		baseURL:  url,
	}
}

// Callback for PubSub messages
func (j *JIRA) Callback(ctx context.Context, m *pubsub.Message) {
	//j.Debug("callback event", zap.ByteString("message", m.Data))
	var msg hangouts.Event
	err := json.Unmarshal(m.Data, &msg)
	if err != nil {
		j.Error("event parse error", zap.Error(err))
		j.hangouts.Send(msg.Space.Name, j.BuildError(errors.Wrap(err, "event parse error"), msg.Message.Thread.Name))
		return
	}
	j.Debug(
		"event received",
		zap.String("type", msg.Type),
		zap.String("space", msg.Space.Name),
		zap.String("user", msg.User.Email),
		zap.String("msg", msg.Message.Text),
	)
	m.Ack()
	regex := regexp.MustCompile("[A-Z]+-[1-9][0-9]*")
	issues := regex.FindAllString(msg.Message.Text, -1)
	for _, issue := range issues {
		j.Debug("handling issue", zap.String("id", issue))
		m, err := j.BuildMessage(issue, msg)
		if err != nil {
			j.Warn("build card failed", zap.Error(err))
			mb, _ := json.Marshal(m)
			j.Debug("message", zap.String("content", string(mb)))
			j.hangouts.Send(msg.Space.Name, j.BuildError(errors.Wrap(err, "build card failed"), msg.Message.Thread.Name))
			return
		}
		err = j.hangouts.Send(msg.Space.Name, m)
		mb, _ := json.Marshal(m)
		j.Debug("message", zap.String("content", string(mb)))
		if err != nil {
			j.Error("send card error", zap.Error(err))
			j.hangouts.Send(msg.Space.Name, j.BuildError(errors.Wrap(err, "send card error"), msg.Message.Thread.Name))
		}
	}
}

// BuildMessage for issue
func (j *JIRA) BuildMessage(issue string, msg hangouts.Event) (*hangouts.Message, error) {
	i, _, err := j.Client.Issue.Get(issue, nil)
	if err != nil {
		j.Debug("jira issue fetch error", zap.String("issue", issue), zap.Error(err))
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
