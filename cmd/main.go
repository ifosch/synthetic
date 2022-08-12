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
	)
	if err := jenkins.Connect(); err != nil {
		log.Fatalf("error connecting to jenkins: %s", err.Error())
	}

	registerChatCommands(chat)
	registerJenkinsCommands(chat, jenkins)
	registerK8sCommands(chat)

	chat.Start()
}

func newRegisterChatCommands(handler *command.Handler) {
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

func newRegisterJenkinsCommands(handler *command.Handler, jenkins *jobcontrol.Jenkins) {
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

func newRegisterK8sCommands(handler *command.Handler) {
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
