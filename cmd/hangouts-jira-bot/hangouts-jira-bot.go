package main

import (
	"context"
	"fmt"
	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/kolide/kit/version"
	"github.com/seibert-media/go-hangouts"
	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/seibert-media/hangouts-jira-bot/pkg/jira"
	"github.com/seibert-media/hangouts-jira-bot/pkg/pubsub"
)

const (
	appName = "Hangouts Jira Bot"
	appKey  = "hangouts-jira-bot"
)

var (
	maxprocs    = flag.Int("maxprocs", runtime.NumCPU(), "max go procs")
	dbg         = flag.Bool("debug", false, "enable debug mode")
	versionInfo = flag.Bool("version", true, "show version info")
	sentryDsn   = flag.String("sentryDsn", "", "sentry dsn key")

	serviceAccount = flag.String("google-service-account", "auth.json", "path to the service account json file")
	projectID      = flag.String("google-project-id", "", "google cloud project id")
	topic          = flag.String("google-pubsub-topic", "", "google cloud pubsub topic")
	subscription   = flag.String("google-pubsub-subscription", "", "google cloud pubsub subscription")

	jiraURL      = flag.String("jira-url", "https://jira.example.com", "jira base url")
	jiraUsername = flag.String("jira-username", "admin", "jira username")
	jiraPassword = flag.String("jira-password", "admin", "jira password")
)

func main() {
	flag.Parse()

	if *versionInfo {
		v := version.Version()
		fmt.Printf("-- //S/M %s --\n", appName)
		fmt.Printf(" - version: %s\n", v.Version)
		fmt.Printf("   branch: \t%s\n", v.Branch)
		fmt.Printf("   revision: \t%s\n", v.Revision)
		fmt.Printf("   build date: \t%s\n", v.BuildDate)
		fmt.Printf("   build user: \t%s\n", v.BuildUser)
		fmt.Printf("   go version: \t%s\n", v.GoVersion)
	}
	runtime.GOMAXPROCS(*maxprocs)

	var zapFields []zapcore.Field
	if !*dbg {
		zapFields = []zapcore.Field{
			zap.String("app", appKey),
			zap.String("version", version.Version().Version),
		}
	}

	logger := log.New(*sentryDsn, *dbg)
	defer logger.Sync()

	ctx := log.WithFields(context.Background(), zapFields...)

	log.From(ctx).Info("preparing")

	ps, err := pubsub.New(ctx, *serviceAccount, *projectID, *topic, *subscription)
	if err != nil {
		log.From(ctx).Fatal("creating pubsub client", zap.Error(err))
	}

	ha, err := hangouts.New(ctx, *serviceAccount)
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
