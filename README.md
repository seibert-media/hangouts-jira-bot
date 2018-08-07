# Hangouts Jira Bot

A simple bot sending cards containing Jira Issues to the respective Hangouts Chat Thread (Jira Server only) using Google Pub/Sub.

## Why

With the release of Hangouts Chat in Feb 2018 we decided to switch from Hipchat to Chat a few months later. Before being able to do so, we had to reach a comparable level of integration inside Chat, the Jira Server integration being one of them.

While there is a working integration for Jira Cloud, we and most of our customers rely on Jira Server, therefor we built this small bot.

**FYI: This project is merely a prototype matching our internal needs until more interested users and customers arise to give us reasons for further extending it.**

## How

This bot is using Google Pub/Sub for subscribing to Hangouts Chat events and acting accordingly.
Messages are then sent via the Hangouts Chat API, which is done using our [own library](https://github.com/seibert-media/go-hangouts).

## Features

While this project mostly got built in a few days, we still got our main features we need covered, with some being planned for the near future.

### Current

* Inviting the Bot to any Hangouts Chat Channel
* Reacting upon mention and sending issue information for found issue key
* Sending multiple tickets if multiple keys are found

### Planned

* Listening on Jira Issue URL's as well
* Listening for lowercase Issue Keys as well
* Dynamic interactions with issues inside Chat (by using the new features in Chat)
* Using OAuth instead of a static Jira login when initially adding the bot to a new room
* Extending this project further and making it easier to deploy through the G Suite Marketplace

## Deployment

We currently run this bot inside our own Kubernetes Cluster on GCP. We might publish the related manifest files in the near future but so far there isn't much to it. Just start the container as a deployment inside your GCP Project.

If you want to have this bot for your own instance, we highly recommend to get in touch with us so we can evaluate your usecase and find the best solution. This bot is a service we provide and maintain so providing it as a useful solution and even hosting it is definitely inside our scope for this.
Do not hesitate to ask us questions via GitHub or by sending inquiries to [google@seibert-media.net](mailto:google@seibert-media.net).
To install this bot yourself you need a G Suite instance as well as a Google Cloud Platform account with billing setup and enabled.

Follow the documentation found at Google to setup your Cloud Platform Project, Google Pub/Sub and the Hangouts Chat API: https://developers.google.com/hangouts/chat/how-tos/pub-sub

Then deploy the [docker image](https://quay.io/repository/seibertmedia/hangouts-jira-bot) and pass the following required parameters as environment variables:

* **GOOGLE_SERVICE_ACCOUNT**: Path to the auth.json file created when setting up your project
* **GOOGLE_PROJECT_ID**: Google Cloud Project ID you setup earlier
* **GOOGLE_PUBSUB_TOPIC**: The PubSub Topic you setup earlier
* **GOOGLE_PUBSUB_SUBSCRIPTION**: The PubSub Subscription you setup earlier

* **JIRA_URL**: URL of your Jira Server Instance
* **JIRA_USERNAME**: Username as which the bot will access Jira to search/get Issues
* **JIRA_PASSWORD**: Password of the above user

Optionally you can set/override the following parameters:

* **SENTRYDSN**: The sentry dsn key for collecting detailed error reports
* **DEBUG**: Whether or not the app should do debug logging (default: false)
* **VERSION**: Whether or not the app should print version info (default: true)

## Usage

Once ready and deployed you can add the bot to any room by mentioning it with the name configured in the setup step. We currently use @jira-server which might change in the near future.

Then once you want to show the ticket(s) shown to Issue Key(s) in your message, just mention the bot in the same message to trigger it. See this example:

```none
Hey team,
I just created a new ticket regarding the readme we got to write for the Hangouts Chat Bot for Jira.
See @jira-server GD-42
```

Sadly until now, the bot always has to be mentioned to react. While there is a feature request for changing this, we don't expect it to happen and honestly think in it's current state this is the better behavior for bots.

## Support

Currently we provide active support only to customers directly working with us. If you got any problems with this project and are not a //SEIBERT/MEDIA GmbH customer, please file your issues here on GitHub and we will try to help you.

## Contributing

This project is currently maintained by a very small subset of the Google team at //SEIBERT/MEDIA.
While rising interest in this integration may lead to more investment on our side, we are always open for users and developers who are open to help with testing and extending this project.
If you are interested in contributing, please file your ideas as issues on GitHub and discuss them with us. PR's are very welcome.

## Attributions

* [Kolide for providing `kit`](https://github.com/kolide/kit)

## License

This project's license is located in the [LICENSE file](LICENSE).