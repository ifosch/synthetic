package main

import (
	"log"
	"os"

	"github.com/ifosch/synthetic-david/pkg/slack"
)

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false
	slack.Start(slackToken, debug)
}
