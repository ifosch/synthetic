package main

import (
	"log"
	"os"

	"github.com/ifosch/synthetic/pkg/jenkins"
	"github.com/ifosch/synthetic/pkg/slack"
	"github.com/ifosch/synthetic/pkg/synthetic"
)

func replyHello(msg synthetic.Message) {
	msg.Reply("hello", msg.Thread())
}

func reactHello(msg synthetic.Message) {
	msg.React("wave")
}

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false

	client := slack.NewChat(slackToken, true, debug)
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/main.replyHello",
			slack.Mentioned(slack.Contains(replyHello, "hello")),
		),
	)
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/main.reactHello",
			slack.NotMentioned(slack.Contains(reactHello, "hello")),
		),
	)

	jenkins.Register(client)

	client.Start()
}
