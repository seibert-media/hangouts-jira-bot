package pubsub

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// PubSub connection
type PubSub struct {
	*pubsub.Client
	*pubsub.Topic
	*pubsub.Subscription
}

// New PubSub connection
func New(ctx context.Context, serviceAccount, projectID, topic, subscription string) (*PubSub, error) {
	ctx = log.WithFields(ctx,
		zap.String("component", "pubsub"),
		zap.String("projectID", projectID),
		zap.String("topic", topic),
		zap.String("subscription", subscription),
	)

	client, err := pubsub.NewClient(ctx, projectID, option.WithServiceAccountFile(serviceAccount))
	if err != nil {
		log.From(ctx).Error("pubsub connection error", zap.Error(err))
		return nil, err
	}
	log.From(ctx).Info("pubsub connected")

	t := client.Topic(topic)
	if t == nil {
		log.From(ctx).Error("topic not found")
		return nil, errors.New("topic not found")
	}
	log.From(ctx).Info("topic connected")

	s := client.Subscription(subscription)
	if s == nil {
		log.From(ctx).Error("subscription not found")
		return nil, errors.New("subscription not found")
	}
	log.From(ctx).Info("subscription connected")

	p := &PubSub{
		Client:       client,
		Topic:        t,
		Subscription: s,
	}
	return p, nil
}
