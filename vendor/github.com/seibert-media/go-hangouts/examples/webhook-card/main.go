package main

import (
	"context"
	"flag"
	"fmt"

	"go.uber.org/zap"

	"github.com/seibert-media/go-hangouts"
	"github.com/seibert-media/golibs/log"
)

var (
	url = flag.String("url", "", "webhook url provided by hangouts chat")
)

func main() {
	flag.Parse()

	fmt.Println(*url)

	// Build Context with Debug Logger
	ctx := log.WithLogger(context.Background(), log.New("", true))

	// Build the card message
	msg := hangouts.NewMessage().WithCard(
		hangouts.NewCard().WithHeader(
			"Example Card",
			"Subtitle",
			"",
			"",
		).WithSection(
			hangouts.NewSection("").WithWidget(
				hangouts.NewWidget().WithTextParagraph("Some Text"),
			),
		).WithSection(
			hangouts.NewSection("").WithWidget(hangouts.NewWidget().
				WithButton(hangouts.NewTextLinkButton(
					"Some Link",
					"https://github.com/seibert-media/go-hangouts"),
				),
			),
		),
	)

	log.From(ctx).Debug("create client", zap.String("url", *url))
	client, err := hangouts.NewWebhookClient(*url)
	if err != nil {
		log.From(ctx).Fatal("create client", zap.Error(err))
	}

	log.From(ctx).Debug("send message", zap.String("url", *url))
	err = client.Send(ctx, "", msg)
	if err != nil {
		log.From(ctx).Fatal("send message", zap.Error(err))
	}

}
