package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// Message is ...
type Message struct {
	event        *slack.MessageEvent
	chat         *Chat
	Completed    bool
	Thread       bool
	Mention      bool
	User         *User
	Conversation *Conversation
	Text         string
}

// ReadMessage ...
func ReadMessage(event *slack.MessageEvent, chat *Chat) (msg *Message, err error) {
	thread := false
	if event.ClientMsgID == "" {
		return &Message{
			event:        event,
			chat:         chat,
			Completed:    false,
			Thread:       thread,
			Mention:      false,
			User:         nil,
			Conversation: nil,
			Text:         "",
		}, nil
	}
	if event.ThreadTimestamp != "" {
		thread = true
	}
	user, err := NewUserFromID(event.User, chat.api)
	if err != nil {
		return nil, err
	}
	conversation, err := NewConversationFromID(event.Channel, chat.api)
	if err != nil {
		return nil, err
	}
	return &Message{
		event:        event,
		chat:         chat,
		Completed:    true,
		Thread:       thread,
		Mention:      strings.Contains(event.Text, fmt.Sprintf("<@%v>", chat.botID)),
		User:         user,
		Conversation: conversation,
		Text:         event.Text,
	}, nil
}

// Reply ...
func (m *Message) Reply(msg string, inThread bool) {
	var message *slack.OutgoingMessage
	if inThread || m.Thread {
		message = m.chat.rtm.NewOutgoingMessage(msg, m.event.Channel, slack.RTMsgOptionTS(m.event.ThreadTimestamp))
	} else if m.chat.defaultReplyInThread {
		message = m.chat.rtm.NewOutgoingMessage(msg, m.event.Channel, slack.RTMsgOptionTS(m.event.Timestamp))
	} else {
		message = m.chat.rtm.NewOutgoingMessage(msg, m.event.Channel)
	}
	m.chat.rtm.SendMessage(message)
}

// React ...
func (m *Message) React(reaction string) {
	m.chat.api.AddReaction(reaction, slack.ItemRef{Channel: m.event.Channel, Timestamp: m.event.Timestamp})
}

// ClearMention ...
func (m *Message) ClearMention() string {
	if !m.Mention {
		return m.Text
	}
	return RemoveWord(m.Text, fmt.Sprintf("<@%v>", m.chat.botID))
}
