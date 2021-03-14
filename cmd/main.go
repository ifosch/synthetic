package main

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

// User is ...
type User struct {
	slackUser *slack.User
	Name      string
}

// NewUserFromID ...
func NewUserFromID(id string, api *slack.Client) (user *User, err error) {
	userInfo, err := api.GetUserInfo(id)
	if err != nil {
		return nil, err
	}
	user = &User{userInfo, fmt.Sprintf("@%v", userInfo.Name)}
	return user, err
}

// Conversation is ...
type Conversation struct {
	slackChannel *slack.Channel
	Name         string
}

// NewConversationFromID ...
func NewConversationFromID(id string, api *slack.Client) (conversation *Conversation, err error) {
	conversationInfo, err := api.GetConversationInfo(id, false)
	if err != nil {
		return nil, err
	}
	conversationName := "DM"
	if conversationInfo.IsChannel {
		conversationName = fmt.Sprintf("#%v", conversationInfo.Name)
		if conversationInfo.Purpose.Value != "" {
			conversationName = conversationInfo.Purpose.Value
		}
	}
	conversation = &Conversation{conversationInfo, conversationName}
	return conversation, nil
}

// Message is ...
type Message struct {
	event        *slack.MessageEvent
	Completed    bool
	Thread       bool
	User         *User
	Conversation *Conversation
	Text         string
}

// ReadMessage ...
func ReadMessage(event *slack.MessageEvent, api *slack.Client) (msg *Message, err error) {
	thread := false
	if event.ClientMsgID == "" {
		return &Message{
			event:        event,
			Completed:    false,
			Thread:       thread,
			User:         nil,
			Conversation: nil,
			Text:         "",
		}, nil
	}
	if event.ThreadTimestamp != "" {
		thread = true
	}
	user, err := NewUserFromID(event.User, api)
	if err != nil {
		return nil, err
	}
	conversation, err := NewConversationFromID(event.Channel, api)
	if err != nil {
		return nil, err
	}
	return &Message{
		event:        event,
		Completed:    true,
		Thread:       thread,
		User:         user,
		Conversation: conversation,
		Text:         event.Text,
	}, nil
}

func main() {
	slackToken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatalf("No SLACK_TOKEN environment variable defined")
	}
	debug := false
	api := slack.New(
		slackToken,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			msg, err := ReadMessage(ev, api)
			if err != nil {
				log.Printf("Error %v processing message %v", err, ev)
				break
			}
			if msg.Completed {
				thread := ""
				if msg.Thread {
					thread = "a thread in "
				}
				log.Printf("Message: '%v' from '%v' in %v'%v'\n", msg.Text, msg.User.Name, thread, msg.Conversation.Name)
			}
		}
	}
}
