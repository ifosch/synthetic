package slack

import (
	"github.com/slack-go/slack"
)

// IClient is an interface for the chat system's client.
type IClient interface {
	GetConversationInfo(string, bool) (*slack.Channel, error)
	GetUserInfo(string) (*slack.User, error)
	NewRTM(...slack.RTMOption) *slack.RTM
	AddReaction(string, slack.ItemRef) error
	RemoveReaction(string, slack.ItemRef) error
}
