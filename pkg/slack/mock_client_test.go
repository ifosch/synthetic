package slack

import (
	"github.com/slack-go/slack"
)

type reactionData struct {
	reaction string
	item     slack.ItemRef
}

type MockClient struct {
	channels         map[string]*slack.Channel
	users            map[string]*slack.User
	reactionsAdded   []reactionData
	reactionsRemoved []reactionData
}

func (c *MockClient) GetConversationInfo(id string, includeLocale bool) (channel *slack.Channel, err error) {
	return c.channels[id], nil
}

func (c *MockClient) GetUserInfo(id string) (*slack.User, error) {
	return c.users[id], nil
}

func (c *MockClient) NewRTM(options ...slack.RTMOption) *slack.RTM {
	return nil
}

func (c *MockClient) AddReaction(reaction string, item slack.ItemRef) error {
	c.reactionsAdded = append(c.reactionsAdded, reactionData{
		reaction: reaction,
		item:     item,
	})
	return nil
}

func (c *MockClient) RemoveReaction(reaction string, item slack.ItemRef) error {
	c.reactionsRemoved = append(c.reactionsRemoved, reactionData{
		reaction: reaction,
		item:     item,
	})
	return nil
}

func (c *MockClient) reset() {
	c.channels = map[string]*slack.Channel{
		"CH00001": &slack.Channel{
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
		"DM00001": &slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "DM00001",
				},
			},
			IsChannel: false,
		},
		"GR00001": &slack.Channel{
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
		"U000001": &slack.User{
			ID:   "U000001",
			Name: "username",
		},
	}
	c.reactionsAdded = []reactionData{}
	c.reactionsRemoved = []reactionData{}
}

func NewMockClient() *MockClient {
	client := &MockClient{}
	client.reset()
	return client
}
