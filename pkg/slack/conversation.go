package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// Conversation is a wrapper over slack-go's Channel object. It
// provides an abstraction layer over channels, group conversations
// and direct chats.
type Conversation struct {
	slackChannel *slack.Channel
	name         string
}

// NewConversationFromID returns a Conversation object wrapping the
// channel, group conversation, or direct chat identified by `id`.
func NewConversationFromID(id string, api IClient) (conversation *Conversation, err error) {
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

// Name returns the name of the conversation.
func (c *Conversation) Name() string {
	return c.name
}
