package main

import (
	"log"
	"os"

	"github.com/ifosch/synthetic-david/pkg/slack"
)

func replyHello(msg *slack.Message) {
	if msg.Text == "hello" {
		msg.Reply("hello", msg.Thread)
	}
}
	}
	msg.Reply("hello", msg.Thread)
}

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false
	client := slack.NewChat(slackToken, true, debug)
	client.RegisterMessageProcessor(replyHello)
	client.Start()
}
