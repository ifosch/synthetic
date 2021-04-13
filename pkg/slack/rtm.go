package slack

import (
	"github.com/slack-go/slack"
)

// IRTM is an interface for the chat system RTM interface.
type IRTM interface {
	ManageConnection()
	NewOutgoingMessage(string, string, ...slack.RTMsgOption) *slack.OutgoingMessage
	SendMessage(*slack.OutgoingMessage)
}
