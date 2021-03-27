package main

import (
	"log"
	"os"
	"strings"

	"github.com/ifosch/synthetic-david/pkg/slack"
)

func replyHello(msg *slack.Message) {
	if msg.Mention && strings.Contains(msg.Text, "hello") {
		msg.Reply("hello", msg.Thread)
	}
}

func reactHello(msg *slack.Message) {
	if !msg.Mention && strings.Contains(msg.Text, "hello") {
		msg.React("wave")
	}
}

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false
	client := slack.NewChat(slackToken, true, debug)
	client.RegisterMessageProcessor(replyHello)
	client.RegisterMessageProcessor(reactHello)
	client.Start()
}
