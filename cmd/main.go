package main

import (
	"log"
	"os"

	"github.com/slack-go/slack"

	jobcontrol "github.com/ifosch/synthetic/pkg/job_control"
	"github.com/ifosch/synthetic/pkg/k8s"
	myslack "github.com/ifosch/synthetic/pkg/slack"
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

	// Initialize dependencies
	api := slack.New(
		slackToken,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	chat := myslack.NewChat(
		api,
		true,
		"",
	)

	jenkins := jobcontrol.NewJenkins(
		os.Getenv("JENKINS_URL"),
		os.Getenv("JENKINS_USER"),
		os.Getenv("JENKINS_PASSWORD"),
		&jobcontrol.JenkinsJobServer{},
	)
	if err := jenkins.Connect(); err != nil {
		log.Fatalf("error connecting to jenkins: %s", err.Error())
	}

	registerChatCommands(chat)
	registerJenkinsCommands(chat, jenkins)
	registerK8sCommands(chat)

	chat.Start()
}

func registerChatCommands(chat *myslack.Chat) {
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/pkg/slack.LogMessage",
			myslack.LogMessage,
		),
	)

	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/main.replyHello",
			myslack.Mentioned(myslack.Contains(replyHello, "hello")),
		),
	)
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/main.reactHello",
			myslack.NotMentioned(myslack.Contains(reactHello, "hello")),
		),
	)
}

func registerJenkinsCommands(chat *myslack.Chat, jenkins *jobcontrol.Jenkins) {
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.List",
			myslack.Exactly(myslack.Mentioned(jenkins.List), "list"),
		),
	)
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.Describe",
			myslack.Mentioned(myslack.Contains(jenkins.Describe, "describe")),
		),
	)
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.Build",
			myslack.Mentioned(myslack.Contains(jenkins.Build, "build")),
		),
	)
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/jenkins.Reload",
			myslack.Mentioned(myslack.Contains(jenkins.Reload, "reload")),
		),
	)
}

func registerK8sCommands(chat *myslack.Chat) {
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/k8s.listClusters",
			myslack.Exactly(myslack.Mentioned(k8s.ListClusters), "list clusters"),
		),
	)
	chat.RegisterMessageProcessor(
		myslack.NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/k8s.listPods",
			myslack.Contains(myslack.Mentioned(k8s.ListPods), "list pods"),
		),
	)
}
