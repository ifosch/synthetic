package slack

import (
	"github.com/slack-go/slack"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

// Message contains all the information about a message the bot was
// notified about.
type Message struct {
	event        *slack.MessageEvent
	chat         *Chat
	Completed    bool
	thread       bool
	mention      bool
	user         *User
	conversation *Conversation
	text         string
}

// Thread is an accessor for Thread.
func (m *Message) Thread() bool {
	return m.thread
}

// Mention is an accessor for Mention.
func (m *Message) Mention() bool {
	return m.mention
}

// User is an accessor for User.
func (m *Message) User() synthetic.User {
	return m.user
}

// Conversation is an accessor for Conversation.
func (m *Message) Conversation() synthetic.Conversation {
	return m.conversation
}

// Text is an accessor for text.
func (m *Message) Text() string {
	return m.text
}

// Reply send the `msg` string as a reply to the message, in a thread
// if `inThread` is true.
func (m *Message) Reply(msg string, inThread bool) {
	var message *slack.OutgoingMessage
	if inThread || m.thread {
		message = m.chat.rtm.NewOutgoingMessage(msg, m.event.Channel, slack.RTMsgOptionTS(m.event.ThreadTimestamp))
	} else if m.chat.defaultReplyInThread {
		message = m.chat.rtm.NewOutgoingMessage(msg, m.event.Channel, slack.RTMsgOptionTS(m.event.Timestamp))
	} else {
		message = m.chat.rtm.NewOutgoingMessage(msg, m.event.Channel)
	}
	m.chat.rtm.SendMessage(message)
}

// React adds the `reaction` reaction to the message.
func (m *Message) React(reaction string) {
	m.chat.api.AddReaction(reaction, slack.ItemRef{Channel: m.event.Channel, Timestamp: m.event.Timestamp})
}

// Unreact removes the `reaction` reaction from the message.
func (m *Message) Unreact(reaction string) {
	m.chat.api.RemoveReaction(reaction, slack.ItemRef{Channel: m.event.Channel, Timestamp: m.event.Timestamp})
}
