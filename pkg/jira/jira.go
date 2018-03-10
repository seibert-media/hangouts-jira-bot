package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/seibert-media/hangouts-jira-bot/pkg/hangouts"

	"cloud.google.com/go/pubsub"
	goJira "github.com/andygrunwald/go-jira"
	"github.com/playnet-public/libs/log"
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
		card, err := j.BuildCard(issue)
		if err != nil {
			j.Warn("build card failed", zap.Error(err))
			return
		}
		err = j.hangouts.SendCard(msg.Space.Name, card)
		if err != nil {
			j.Error("send card error", zap.Error(err))
		}
	}
}

// BuildCard for issue
func (j *JIRA) BuildCard(issue string) (string, error) {
	i, _, err := j.Client.Issue.Get(issue, nil)
	if err != nil {
		j.Debug("jira issue fetch error", zap.String("issue", issue), zap.Error(err))
		return "", err
	}
	/*card := hangouts.Card{
		Header: hangouts.Header{
			Title:      issue,
			Subtitle:   i.Fields.Summary,
			ImageURL:   "https://cdn6.aptoide.com/imgs/9/9/b/99b698eae5433cc15b23862f4a305a37_icon.png?w=240",
			ImageStyle: "AVATAR",
		},
		Sections: []hangouts.Section{
			{
				Widgets: []hangouts.Widget{
					{
						KeyValue: hangouts.KeyValue{
							TopLabel: "Status",
							Content:  i.Fields.Status.Name,
						},
					},
					{
						TextParagraph: hangouts.TextParagraph{
							Text: i.Fields.Description,
						},
					},
				},
			},
			{
				Widgets: []hangouts.Widget{
					hangouts.ActionWidget{
						Buttons: []hangouts.Button{
							hangouts.Button{
								TextButton: hangouts.TextButton{
									Text: "Open Ticket",
									OnClick: hangouts.OnClick{
										OpenLink: hangouts.OpenLink{
											URL: i.Self,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Actions: []hangouts.Action{
			{
				Label: "Open Ticket",
				OnClick: hangouts.OnClick{
					OpenLink: hangouts.OpenLink{
						URL: i.Self,
					},
				},
			},
		},
	}*/
	card := fmt.Sprintf(`{ "cards": [
		{
		"header": {
			"title": "%s",
			"subtitle": "%s",
			"imageUrl": "https://storage.googleapis.com/bot-icons/jira-app-icon.png",
			"imageStyle": "AVATAR"
		},
		"sections": [{
				"widgets": [{
						"keyValue": {
							"topLabel": "Status",
							"content": "%s"
						}
					},
					{
						"textParagraph": {
							"text": "%s"
						}
					}
				]
			},
			{
				"widgets": [{
					"buttons": [{
						"textButton": {
							"text": "Open Ticket",
							"onClick": {
								"openLink": {
									"url": "%s/browse/%s"
								}
							}
						}
					}]
				}]
			}
		]
	}]}`, issue, i.Fields.Summary, i.Fields.Status.Name, i.Fields.Description, j.baseURL, issue)
	return card, nil
}
