package slack

import (
	"github.com/slack-go/slack"
)

type reactionData struct {
	reaction string
	item     slack.ItemRef
}

// MockClient is a mocking client for testing.
type MockClient struct {
	channels         map[string]*slack.Channel
	users            map[string]*slack.User
	reactionsAdded   []reactionData
	reactionsRemoved []reactionData
}

// GetConversationInfo returns the channel information for `id`.
func (c *MockClient) GetConversationInfo(id string, includeLocale bool) (channel *slack.Channel, err error) {
	return c.channels[id], nil
}

// GetUserInfo returns the user information for `id`.
func (c *MockClient) GetUserInfo(id string) (*slack.User, error) {
	return c.users[id], nil
}

// NewRTM returns a null Slack RTM.
func (c *MockClient) NewRTM(options ...slack.RTMOption) *slack.RTM {
	return nil
}

// AddReaction registers `reaction` on `item` for validation.
func (c *MockClient) AddReaction(reaction string, item slack.ItemRef) error {
	c.reactionsAdded = append(c.reactionsAdded, reactionData{
		reaction: reaction,
		item:     item,
	})
	return nil
}

// RemoveReaction registers `reaction` removal on `item` for validation.
func (c *MockClient) RemoveReaction(reaction string, item slack.ItemRef) error {
	c.reactionsRemoved = append(c.reactionsRemoved, reactionData{
		reaction: reaction,
		item:     item,
	})
	return nil
}

func (c *MockClient) reset() {
	c.channels = map[string]*slack.Channel{
		"CH00001": {
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "CH00001",
				},
				Name: "test",
				Purpose: slack.Purpose{
					Value: "",
				},
			},
			IsChannel: true,
		},
		"DM00001": {
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "DM00001",
				},
			},
			IsChannel: false,
		},
		"GR00001": {
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "GR00001",
				},
				Name: "",
				Purpose: slack.Purpose{
					Value: "Group messaging with: @some @users",
				},
			},
			IsChannel: true,
		},
	}
	c.users = map[string]*slack.User{
		"U000001": {
			ID:   "U000001",
			Name: "username",
		},
	}
	c.reactionsAdded = []reactionData{}
	c.reactionsRemoved = []reactionData{}
}

// NewMockClient creates a new MockClient.
func NewMockClient() *MockClient {
	client := &MockClient{}
	client.reset()
	return client
}

// MockRTM is a mocking RTM.
type MockRTM struct {
	messagesSent []*slack.OutgoingMessage
}

// ManageConnection fakes the real Slack RTM connection manager.
func (rtm *MockRTM) ManageConnection() {}

// NewOutgoingMessage creates a fake message object to send.
func (rtm *MockRTM) NewOutgoingMessage(text string, channelID string, options ...slack.RTMsgOption) *slack.OutgoingMessage {
	msg := &slack.OutgoingMessage{
		Channel: channelID,
		Text:    text,
	}
	for _, option := range options {
		option(msg)
	}
	return msg
}

// SendMessage fakes sending a message.
func (rtm *MockRTM) SendMessage(msg *slack.OutgoingMessage) {
	rtm.messagesSent = append(rtm.messagesSent, msg)
}

func (rtm *MockRTM) reset() {
	rtm.messagesSent = []*slack.OutgoingMessage{}
}

// NewMockRTM creates a new MockRTM.
func NewMockRTM() *MockRTM {
	rtm := &MockRTM{}
	rtm.reset()
	return rtm
}
