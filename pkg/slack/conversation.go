package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

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
