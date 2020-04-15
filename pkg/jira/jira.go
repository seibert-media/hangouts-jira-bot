package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"cloud.google.com/go/pubsub"
	goJira "github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	"google.golang.org/api/chat/v1"
)

// JIRA connection
type JIRA struct {
	*goJira.Client
	chat    *chat.SpacesMessagesService
	baseURL string
}

// New JIRA connection
func New(ctx context.Context, url, username, password string, chat *chat.Service) *JIRA {
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
		Client:  client,
		chat:    chat.Spaces.Messages,
		baseURL: url,
	}
}

// Callback for PubSub messages
func (j *JIRA) Callback(ctx context.Context, m *pubsub.Message) {
	m.Ack()

	var event *chat.DeprecatedEvent
	err := json.Unmarshal(m.Data, &event)
	if err != nil {
		log.From(ctx).Error("parsing event", zap.Error(err))
		j.handleDetailedError(ctx, nil, err, "parsing message")
		return
	}

	if event == nil || event.Message == nil {
		log.From(ctx).Warn("skipping empty message", zap.String("raw", string(m.Data)))
		return
	}

	ctx = log.WithFields(ctx,
		zap.String("type", event.Type),
		zap.String("space", event.Space.Name),
		zap.String("user", event.User.Name),
	)

	msg := event.Message
	log.From(ctx).Debug(
		"receiving event",
		zap.String("msg", msg.Text),
	)

	regex := regexp.MustCompile("[A-Za-z]+-[1-9][0-9]*")
	issues := regex.FindAllString(msg.Text, -1)
	for _, issue := range issues {
		ctx := log.WithFields(ctx, zap.String("id", issue))
		log.From(ctx).Debug("handling issue")

		m, err := j.BuildMessage(ctx, issue, msg)
		if err != nil {
			log.From(ctx).Warn("building message", zap.Error(err))
			j.handleDetailedError(ctx, msg, err, "building card failed")
			continue
		}

		if _, err := j.chat.Create(msg.Space.Name, m).Context(ctx).Do(); err != nil {
			log.From(ctx).Error("sending message", zap.Error(err))
			j.handleDetailedError(ctx, msg, err, "sending card failed")
			continue
		}
	}
}

// BuildMessage for issue
func (j *JIRA) BuildMessage(ctx context.Context, issue string, msg *chat.Message) (*chat.Message, error) {
	i, _, err := j.Client.Issue.Get(issue, nil)
	if err != nil {
		log.From(ctx).Debug("fetching issue", zap.String("issue", issue), zap.Error(err))
		return nil, err
	}

	if i == nil || i.Fields == nil {
		log.From(ctx).Error("fetching issue", zap.Error(errNotFound))
		return nil, errors.Wrap(errNotFound, issue)
	}

	assignee := "Unassigned"
	if i.Fields.Assignee != nil {
		assignee = i.Fields.Assignee.Name
	}

	return &chat.Message{
		Thread: msg.Thread,
		Cards: []*chat.Card{
			{
				Header: &chat.CardHeader{
					Title: fmt.Sprintf("%s: %s", issue, i.Fields.Summary),
				},
				Sections: []*chat.Section{
					{
						Widgets: []*chat.WidgetMarkup{
							{
								KeyValue: &chat.KeyValue{
									TopLabel:         "Status",
									Content:          i.Fields.Status.Name,
									ContentMultiline: false,
								},
							},
							{
								KeyValue: &chat.KeyValue{
									TopLabel:         "Assignee",
									Content:          assignee,
									ContentMultiline: false,
								},
							},
						},
					},
					{
						Widgets: []*chat.WidgetMarkup{
							{
								Buttons: []*chat.Button{
									{
										TextButton: &chat.TextButton{
											Text: "Open Issue",
											OnClick: &chat.OnClick{
												OpenLink: &chat.OpenLink{Url: fmt.Sprintf("%s/browse/%s", j.baseURL, issue)},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

var errNotFound = errors.New("not found")

func (j *JIRA) handleDetailedError(ctx context.Context, msg *chat.Message, err error, details string) {
	var mb []byte
	if msg != nil {
		mb, _ = json.Marshal(msg)
		log.From(ctx).Debug("message", zap.String("content", string(mb)))
	}
	j.chat.Create(msg.Space.Name, j.buildErrorMessage(errors.Wrap(err, details), msg.Thread.Name))
}

func (j *JIRA) buildErrorMessage(e error, thread string) *chat.Message {
	return &chat.Message{
		Thread: &chat.Thread{
			Name: thread,
		},
		Text: fmt.Sprintf("Error: %s", e.Error()),
	}
}
