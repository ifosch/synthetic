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
	event slack.MessageEvent
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
			user, err := NewUserFromID(ev.User, api)
			if err != nil {
				log.Printf("Error trying to get user info for '%v'", ev.User)
			}
			conversation, err := NewConversationFromID(ev.Channel, api)
			if err != nil {
				log.Printf("Error trying to get channel info for '%v'", ev.Channel)
			}
			log.Printf("Message: '%v' from '%v' in '%v'\n", ev.Text, user.Name, conversation.Name)
		}
	}
}
