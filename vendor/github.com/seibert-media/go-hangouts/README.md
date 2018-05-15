# Google Hangouts Chat Library for Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/go-hangouts)](https://goreportcard.com/report/github.com/seibert-media/go-hangouts)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seibert-media/go-hangouts)

Shortly after Hangouts Chat became publicly available, we built [seibert-media/hangouts-jira-bot](https://github.com/seibert-media/hangouts-jira-bot) to try out how the bot connectivity works.

While doing this, one of our first steps was to build this mapping to easily create messages.
As we are now aiming to try out building more bots and integrating Chat into our tool landscape it became reasonable to offer this as a standalone library which is well documented.

In addition to this, we added further helpers to easily build messages available in [builder.go](builder.go).

## Requirements

All dependencies are managed using `dep` and checked into git.
To update those dependencies use `dep ensure --update` after you added new imports to your code.

For building or running tests a onetime execution of `make deps` might be required.

## Usage

For now, this package only offers sending messages. To receive messages and events, we suggest using PubSub. This implementation might be added here in the future but is out of scope for now.
To learn how to do this, check the official docs or give our [jira bot](https://github.com/seibert-media/hangouts-jira-bot) a look.

This package offers both a default connection using Google Service Accounts and sending messages via channel webhooks.

The entire connection Part relies on [github.com/playnet-public/libs/log](https://github.com/playnet-public/libs/log) for logging. If you are not using Zap and are an external user, you could simply use `log.NewNop()` when initializing the connection to get an empty no-op logger.

### Service Account

To send messages using a service account, use `New(...)` which is using the Google DefaultClient for taking credential as described [here](https://cloud.google.com/docs/authentication/production#auth-cloud-implicit-go).

To send messages onto a webhook use `NewWebhookClient(...)`.

Afterwards you can simply use `Send(...)` to post your messages.
Check out code comments for further information on how to use this.

## Examples

