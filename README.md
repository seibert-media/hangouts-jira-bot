# Hangouts Jira Bot

[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/hangouts-jira-bot)](https://goreportcard.com/report/github.com/seibert-media/hangouts-jira-bot)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/72cc701000034d60a0f5caeace0878de)](https://www.codacy.com/app/seibert-media/hangouts-jira-bot?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=seibert-media/hangouts-jira-bot&amp;utm_campaign=Badge_Grade)
[![Docker Image](https://img.shields.io/badge/image-quay.io-success "Docker Image")](https://quay.io/repository/seibertmedia/hangouts-jira-bot)
[![GitHub license](https://img.shields.io/badge/license-AGPL-blue.svg)](https://raw.githubusercontent.com/seibert-media/hangouts-jira-bot/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seibert-media/hangouts-jira-bot)

A simple bot sending cards containing Jira Issues to the respective Hangouts Chat Thread (Jira Server only) using Google Pub/Sub.

## Why

With the release of Hangouts Chat in Feb 2018 we decided to switch from Hipchat to Chat a few months later. Before being able to do so, we had to reach a comparable level of integration inside Chat, the Jira Server integration being one of them.

While there is a working integration for Jira Cloud, we and most of our customers rely on Jira Server, therefor we built this small bot.

**FYI: This project is merely a prototype matching our internal needs until more interested users and customers arise to give us reasons for further extending it.**

## How

This bot is using Google Pub/Sub for subscribing to Hangouts Chat events and acting accordingly.
Responses are sent via the Hangouts Chat API.

## Features

While this project was mostly built in a few days, we still got our main features we need covered, with some being planned for the near future.

### Current

* Inviting the Bot to any Hangouts Chat Channel or Direct Message.
* Reacting upon mention and sending issue information for found issue keys.
* Sending multiple tickets if multiple keys are found.

### Planned

* Listening on Jira Issue URL's as well.
* Creating issues through Chat messages.
* Dynamic interactions with issues inside Chat (by using the new features in Chat).
* Using OAuth instead of a static Jira login when initially adding the bot to a new room.
* Extending this project further and making it easier to deploy through the G Suite Marketplace.

## Deployment

If you want to have this bot for your own instance, we highly recommend to get in touch with us so we can evaluate your usecase and find the best solution. This bot is a service we provide and maintain so providing it as a useful solution and even hosting it is definitely inside our scope for this.
Do not hesitate to ask us questions via GitHub or by sending inquiries to [google@seibert-media.net](mailto:google@seibert-media.net).

We currently run this bot inside our own Kubernetes Cluster on GCP.
The builds and deployments are done by [Google Cloud Build](http://cloud.google.com/cloud-build).

The project contains two cloudbuild configurations.

- `cloudbuild.yaml` is the default config that could be used by consumers to deploy their own bot.
- `cloudbuild_publish.yaml` also publishes the docker image to quay.io for public use and is being run in our own pipeline.

The `cloudbuild.yaml` contains several substitution variables to customize the deployment.
Internally we run three different deployments based on the `dev`, `staging` and `master` branch, each deploying to a different cluster with a different Chat config.

### Project Setup

To deploy the bot, we recommend creating a new Google Cloud Project for it. Make sure to enable billing for the project.
Inside this project, [first enable the Google Cloud Pub/Sub API, then create a topic and subscription](https://console.cloud.google.com/cloudpubsub/topic) both named `hangouts-jira-bot`.

Grant the `Pub/Sub Publisher` role to `chat-api-push@system.gserviceaccount.com` so Chat is allowed to publish messages to the topic.

Afterwards, enable the Hangouts Chat API and [start configuring your bot](https://console.cloud.google.com/apis/api/chat.googleapis.com/hangouts-chat).
You can follow this config:

- **Bot name:** jira
- **Avatar URL:** We recommend to upload a Jira Icon to Google Cloud Storage and insert the public URL here.
- **Description:** Chat bot for Jira Server
- **Functionality:** Check both boxes to allow people using the bot in rooms and direct messages.
- **Connection Settings:** Select `Cloud Pub/Sub` and enter your topic into the textbox. The topic looks like this `projects/your-project-name/topics/hangouts-jira-bot`.
- **Permissions:** The default should be to grant access to everyone in your organization. You can change this if required.

Afterwards, you should [create a new ServiceAccount for the Bot to use](https://console.cloud.google.com/iam-admin/serviceaccounts/create).
Make sure to create a json key and store it securely.
Once created, grant the `Pub/Sub Subscriber` role to the ServiceAccount.

### Kubernetes Deployment

First, create a namespace for your bot (e.g `kubectl create ns hangouts-jira-bot`).

The bot requires certain secrets to run.
To provide those, you can refer to [our own secret template](./kubernetes-manifests/templates/secret.yaml).
This file takes the respective secrets from our internal store and adds them to Kubernetes.
You need to create a respective secret in your namespace before continuing.

Afterwards you can continue with the actual deployment.
The current recomendation is, to fork this repository and setup your own Cloud Build Trigger.
This way, you can utilize the work already done at automizing the deployment and easily update the repository when needed.

To do so, create a fork on GitHub.
Then [link your newly created repository with Cloud Build](https://console.cloud.google.com/cloud-build/triggers/connect).

Once this is done, you can create a trigger.
Make sure to select `Cloud Build configuration file` in the `Build Configuration` section and type `cloudbuild.yaml` into the textbox.
**Do not use the publish yaml, or your build will fail.**

Then add the substitution variables:

- **_CLUSTER_LOCATION:** The location of the Kubernetes cluster you want to deploy to (e.g. `europe-west3-c`).
- **_CLUSTER_NAME:** The name of your Kubernetes cluster (e.g. `cluster-1`).
- **_PUBSUB_PROJECT:** The project of your Pub/Sub Topic (e.g. `hangouts-jira-bot`).
- **_PUBSUB_TOPIC:** The name of your Pub/Sub Topic (e.g. `hangouts-jira-bot`).
- **_PUBSUB_SUBSCRIPTION:** The name of your Pub/Sub Subscription (e.g. `hangouts-jira-bot`).
- **_REGISTRY:** The hostname of your prefered Google Container Registry (e.g. `eu.gcr.io`).

Once you saved, you can manually trigger your first build. Everything should run smoothly. If not, please get in touch. We're here to help.

## Usage

Once ready and deployed you can add the bot to any room by mentioning it with the name configured in the setup step. We currently use @jira.

```none
Hey team,
I just created a new ticket regarding the readme we got to write for the Hangouts Chat Bot for Jira.
See @jira-server GD-42
```

Sadly until now, the bot always has to be mentioned to react. While there is a feature request for changing this, we don't expect it to happen and honestly think in it's current state this is the better behavior for bots.

## Support

If you need help with this project, feel free to raise issues here on GitHub or get in touch with our team at [google@seibert-media.net](mailto:google@seibert-media.net).

## Contributing

This project is currently maintained by a very small subset of the Google Partner team at //SEIBERT/MEDIA.
While rising interest in this integration may lead to more investment on our side, we are always open for users and developers who are open to help with testing and extending this project.
If you are interested in contributing, please file your ideas as issues on GitHub and discuss them with us. PR's are very welcome.

## License

This project's license is located in the [LICENSE file](LICENSE).