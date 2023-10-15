package main

import (
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"

	"github.com/ifosch/synthetic/pkg/command"
	jobcontrol "github.com/ifosch/synthetic/pkg/job_control"
	"github.com/ifosch/synthetic/pkg/k8s"
	myslack "github.com/ifosch/synthetic/pkg/slack"
)

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

	go chat.Start()
	cHandler := command.NewHandler()
	registerChatCommands(cHandler)
	registerJenkinsCommands(cHandler, jenkins)
	registerK8sCommands(cHandler)

	// Blocks until chat.MessageChannel is closed
	cHandler.EventLoop(chat.MessageChannel)
}

func registerChatCommands(handler *command.Handler) {
	var err error
	// LogMessage is a message processor to log the message received.
	err = handler.Register(
		"main.LogMessage",
		func(c *command.Command) {
			msg := c.Message()
			thread := ""
			if msg.Thread() {
				thread = "a thread in "
			}
			log.Printf(
				"Message: '%v' from '%v' in %v'%v'\n",
				msg.Text(),
				msg.User().Name(),
				thread,
				msg.Conversation().Name(),
			)
		},
	)
	if err != nil {
		panic(err)
	}
	err = handler.Register(
		"main.replyHello",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Mention() && strings.Contains(msg.Text(), "hello") {
				msg.Reply("hello", msg.Thread())
			}
		},
	)
	if err != nil {
		panic(err)
	}
	err = handler.Register(
		"main.reactHello",
		func(c *command.Command) {
			msg := c.Message()
			if !msg.Mention() && strings.Contains(msg.Text(), "hello") {
				msg.React("wave")
			}
		},
	)
	if err != nil {
		panic(err)
	}
}

func registerJenkinsCommands(handler *command.Handler, jenkins *jobcontrol.Jenkins) {
	var err error
	err = handler.Register(
		"jenkins.List",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Text() == "list" && msg.Mention() {
				jenkins.List(msg)
			}
		},
	)
	if err != nil {
		panic(err)
	}
	err = handler.Register(
		"jenkins.Describe",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Mention() && strings.Contains(msg.Text(), "describe") {
				jenkins.Describe(msg)
			}
		},
	)
	if err != nil {
		panic(err)
	}
	err = handler.Register(
		"jenkins.Build",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Mention() && strings.Contains(msg.Text(), "build") {
				jenkins.Build(msg)
			}
		},
	)
	if err != nil {
		panic(err)
	}
	err = handler.Register(
		"jenkins.Reload",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Mention() && strings.Contains(msg.Text(), "reload") {
				jenkins.Reload(msg)
			}
		},
	)
	if err != nil {
		panic(err)
	}
}

func registerK8sCommands(handler *command.Handler) {
	var err error
	err = handler.Register(
		"k8s.listClusters",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Text() == "list clusters" && msg.Mention() {
				k8s.ListClusters(msg)
			}
		},
	)
	if err != nil {
		panic(err)
	}
	err = handler.Register(
		"k8s.listPods",
		func(c *command.Command) {
			msg := c.Message()
			if msg.Mention() && strings.Contains(msg.Text(), "list pods") {
				k8s.ListPods(msg)
			}
		},
	)
	if err != nil {
		panic(err)
	}
}
