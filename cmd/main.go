package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/seibert-media/hangouts-jira-bot/pkg/hangouts"

	flag "github.com/bborbe/flagenv"

	"github.com/kolide/kit/version"
	"github.com/playnet-public/libs/log"
	"github.com/seibert-media/hangouts-jira-bot/pkg/jira"
	"github.com/seibert-media/hangouts-jira-bot/pkg/pubsub"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	appName = "hangouts-jira-bot"
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
		fmt.Printf("-- //S/M %s --\n", appName)
		version.PrintFull()
	}
	runtime.GOMAXPROCS(*maxprocs)

	var zapFields []zapcore.Field
	if !*dbg {
		zapFields = []zapcore.Field{
			zap.String("app", appKey),
			zap.String("version", version.Version().Version),
		}
	}

	log := log.New(appKey, *sentryDsn, *dbg).WithFields(zapFields...)
	defer log.Sync()
	log.Info("preparing")

	fmt.Println("main.go")

	err := do(log)
	if err != nil {
		log.Error("application error", zap.Error(err))
		panic(err)
	}
}

func do(log *log.Logger) error {
	ctx := context.Background()
	ps, err := pubsub.New(ctx, log, *serviceAccount, *projectID, *topic, *subscription)
	ha, err := hangouts.New(ctx, log, *serviceAccount)
	if err != nil {
		return err
	}
	ji := jira.New(log, *jiraURL, *jiraUsername, *jiraPassword, ha)

	ps.Info("listening")
	err = ps.Receive(ctx, ji.Callback)
	if err != nil {
		log.Error("receive error", zap.Error(err))
	}
	return err
}
