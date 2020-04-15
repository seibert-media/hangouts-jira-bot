package main

import (
	"context"
	"runtime"

	"github.com/seibert-media/hangouts-jira-bot/pkg/jira"
	"github.com/seibert-media/hangouts-jira-bot/pkg/pubsub"

	flag "github.com/bborbe/flagenv"
	"github.com/seibert-media/go-hangouts"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	var zapFields []zapcore.Field
	if !*dbg {
		zapFields = []zapcore.Field{
			zap.String("app", appKey),
		}
	}

	logger, err := log.New(*sentryDSN, *dbg)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ctx := log.WithLogger(context.Background(), logger)
	ctx = log.WithFields(ctx, zapFields...)

	log.From(ctx).Info("preparing")

	ps, err := pubsub.New(ctx, *projectID, *topic, *subscription)
	if err != nil {
		log.From(ctx).Fatal("creating pubsub client", zap.Error(err))
	}

	ha, err := hangouts.New(ctx, "")
	if err != nil {
		log.From(ctx).Fatal("creating hangouts client", zap.Error(err))
	}

	ji := jira.New(ctx, *jiraURL, *jiraUsername, *jiraPassword, ha)

	log.From(ctx).Info("listening on subscription")
	err = ps.Receive(ctx, ji.Callback)
	if err != nil {
		log.From(ctx).Error("receiving", zap.Error(err))
	}

	log.From(ctx).Info("finished", zap.Error(err))
}
