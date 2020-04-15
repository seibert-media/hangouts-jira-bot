package main

import (
	"context"
	"runtime"

	"github.com/seibert-media/hangouts-jira-bot/pkg/jira"
	"github.com/seibert-media/hangouts-jira-bot/pkg/pubsub"
	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"

	flag "github.com/bborbe/flagenv"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2/google"
)

const appKey = "hangouts-jira-bot"

var (
	maxprocs    = flag.Int("maxprocs", runtime.NumCPU(), "max go procs")
	dbg         = flag.Bool("debug", false, "enable debug mode")
	versionInfo = flag.Bool("version", true, "show version info")
	sentryDSN   = flag.String("sentryDsn", "", "sentry dsn key")

	projectID    = flag.String("google-project-id", "", "google cloud project id")
	topic        = flag.String("google-pubsub-topic", "", "google cloud pubsub topic")
	subscription = flag.String("google-pubsub-subscription", "", "google cloud pubsub subscription")

	jiraURL      = flag.String("jira-url", "https://jira.example.com", "jira base url")
	jiraUsername = flag.String("jira-username", "admin", "jira username")
	jiraPassword = flag.String("jira-password", "admin", "jira password")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*maxprocs)

	logger, err := log.New(*sentryDSN, *dbg)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if *dbg {
		logger.SetLevel(zapcore.DebugLevel)
	}

	ctx := log.WithFields(
		log.WithLogger(context.Background(), logger),
		zap.String("app", appKey),
	)

	log.From(ctx).Info("preparing")

	ps, err := pubsub.New(ctx, *projectID, *topic, *subscription)
	if err != nil {
		log.From(ctx).Fatal("creating pubsub client", zap.Error(err))
	}

	chatClient, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/chat.bot")
	if err != nil {
		log.From(ctx).Fatal("creating chat http client", zap.Error(err))
	}

	chat, err := chat.NewService(ctx, option.WithHTTPClient(chatClient))
	if err != nil {
		log.From(ctx).Fatal("creating chat client", zap.Error(err))
	}

	jira := jira.New(ctx, *jiraURL, *jiraUsername, *jiraPassword, chat)

	log.From(ctx).Info("listening on subscription")
	err = ps.Receive(ctx, jira.Callback)
	if err != nil {
		log.From(ctx).Error("receiving", zap.Error(err))
	}

	log.From(ctx).Info("finished", zap.Error(err))
}
