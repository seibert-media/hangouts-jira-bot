package pubsub

import (
	"context"
	"errors"

	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
	"github.com/playnet-public/libs/log"
	"go.uber.org/zap"
)

// PubSub connection
type PubSub struct {
	*log.Logger
	*pubsub.Client
	*pubsub.Topic
	*pubsub.Subscription
}

// New PubSub connection
func New(ctx context.Context, log *log.Logger, serviceAccount, projectID, topic, subscription string) (*PubSub, error) {
	log = log.WithFields(
		zap.String("component", "pubsub"),
		zap.String("projectID", projectID),
		zap.String("topic", topic),
		zap.String("subscription", subscription),
	)

	client, err := pubsub.NewClient(ctx, projectID, option.WithServiceAccountFile(serviceAccount))
	if err != nil {
		log.Error("pubsub connection error", zap.Error(err))
		return nil, err
	}
	log.Info("pubsub connected")

	t := client.Topic(topic)
	if t == nil {
		log.Error("topic not found")
		return nil, errors.New("topic not found")
	}
	log.Info("topic connected")

	s := client.Subscription(subscription)
	if s == nil {
		log.Error("subscription not found")
		return nil, errors.New("subscription not found")
	}
	log.Info("subscription connected")

	p := &PubSub{
		Logger:       log,
		Client:       client,
		Topic:        t,
		Subscription: s,
	}
	return p, nil
}
