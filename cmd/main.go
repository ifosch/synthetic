package main

import (
	"log"
	"os"

	"github.com/ifosch/synthetic/pkg/jenkins"
	"github.com/ifosch/synthetic/pkg/slack"
)

func replyHello(msg *slack.Message) {
	msg.Reply("hello", msg.Thread)
}

func reactHello(msg *slack.Message) {
	msg.React("wave")
}

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false

	client := slack.NewChat(slackToken, true, debug)
	client.RegisterMessageProcessor(slack.Mentioned(slack.Contains(replyHello, "hello")))
	client.RegisterMessageProcessor(slack.NotMentioned(slack.Contains(reactHello, "hello")))

	jenkins.Register(client)

	client.Start()
}
