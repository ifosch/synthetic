package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ifosch/synthetic-david/pkg/slack"
)

func printMessage(msg *slack.Message) {
	thread := ""
	if msg.Thread {
		thread = "a thread in "
	}
	fmt.Printf("Message: '%v' from '%v' in %v'%v'\n", msg.Text, msg.User.Name, thread, msg.Conversation.Name)
}

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false
	client := slack.NewChat(slackToken, debug)
	client.RegisterMessageProcessor(printMessage)
	client.Start()
}
