package hangouts

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/playnet-public/libs/log"
	"go.uber.org/zap"
	googleAuth "golang.org/x/oauth2/google"
)

// Hangouts handler
type Hangouts struct {
	*log.Logger
	*http.Client
}

// New Hangouts client
func New(ctx context.Context, log *log.Logger, serviceAccount string) (*Hangouts, error) {
	log = log.WithFields(
		zap.String("component", "hangouts"),
	)
	httpClient, err := googleAuth.DefaultClient(ctx, "https://www.googleapis.com/auth/chat.bot")
	if err != nil {
		return nil, err
	}
	return &Hangouts{
		Logger: log,
		Client: httpClient,
	}, nil
}

// SendCard to Hangouts
func (h *Hangouts) SendCard(space string, card string) error {
	/*msg := Message{
		Cards: []Card{
			card,
		},
	}*/
	url := fmt.Sprintf("https://chat.googleapis.com/v1/%s/messages", space)
	//data, err := json.Marshal(msg)
	//if err != nil {
	//	return err
	//}
	resp, err := h.Client.Post(url, "application/json", bytes.NewBuffer([]byte(card)))
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(resp.Body)
		h.Error("post card error", zap.String("status", resp.Status), zap.ByteString("body", body))
		return err
	}
	return nil
}
