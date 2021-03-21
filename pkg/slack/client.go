package slack

import (
	"github.com/slack-go/slack"
)

// IClient ...
type IClient interface {
	GetConversationInfo(string, bool) (*slack.Channel, error)
}
