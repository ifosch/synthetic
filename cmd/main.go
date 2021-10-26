package main

import (
	"log"
	"os"

	jobcontrol "github.com/ifosch/synthetic/pkg/job_control"
	"github.com/ifosch/synthetic/pkg/k8s"
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

	jenkins := jobcontrol.NewJenkins(
		os.Getenv("JENKINS_URL"),
		os.Getenv("JENKINS_USER"),
		os.Getenv("JENKINS_PASSWORD"),
	)
	if err := jenkins.Connect(); err != nil {
		log.Fatalf("error connecting to jenkins: %s", err.Error())
	}
	registerJenkinsCommands(client, jenkins)

	registerK8sCommands(client)

	client.Start()
}

func registerJenkinsCommands(client *slack.Chat, jenkins *jobcontrol.Jenkins) {
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.List",
			slack.Exactly(slack.Mentioned(jenkins.List), "list"),
		),
	)
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.Describe",
			slack.Mentioned(slack.Contains(jenkins.Describe, "describe")),
		),
	)
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.Build",
			slack.Mentioned(slack.Contains(jenkins.Build, "build")),
		),
	)
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.Reload",
			slack.Mentioned(slack.Contains(jenkins.Reload, "reload")),
		),
	)
}

func registerK8sCommands(client *slack.Chat) {
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/k8s.listClusters",
			slack.Exactly(slack.Mentioned(k8s.ListClusters), "list clusters"),
		),
	)
	client.RegisterMessageProcessor(
		slack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/k8s.listPods",
			slack.Contains(slack.Mentioned(k8s.ListPods), "list pods"),
		),
	)
}
