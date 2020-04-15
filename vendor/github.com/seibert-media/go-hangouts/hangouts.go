package hangouts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	googleAuth "golang.org/x/oauth2/google"
)

// Hangouts handler
type Hangouts struct {
	*http.Client
	URL string
}

// New Hangouts client
func New(ctx context.Context, serviceAccount string) (*Hangouts, error) {
	ctx = log.WithFields(ctx, zap.String("component", "hangouts"))
	httpClient, err := googleAuth.DefaultClient(ctx, "https://www.googleapis.com/auth/chat.bot")
	if err != nil {
		return nil, err
	}
	return &Hangouts{
		Client: httpClient,
	}, nil
}

// NewWebhookClient for Hangouts
func NewWebhookClient(url string) (*Hangouts, error) {
	return &Hangouts{
		Client: &http.Client{},
		URL:    url,
	}, nil
}

// Send Message to Hangouts
// space can be left empty if using webhooks as it is used to identify the channel messages are being sent to
func (h *Hangouts) Send(ctx context.Context, space string, msg *Message) error {
	url := h.URL
	if url == "" {
		url = fmt.Sprintf("https://chat.googleapis.com/v1/%s/messages", space)
	}
	ctx = log.WithFields(ctx, zap.String("url", url))
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := h.Client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.From(ctx).Error("post message error", zap.String("status", resp.Status), zap.ByteString("body", body))
		return fmt.Errorf("post message error: %s", body)
	}
	return nil
}
