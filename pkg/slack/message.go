package slack

import (
	"github.com/slack-go/slack"
)

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
