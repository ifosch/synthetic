package slack

import (
	"testing"

	"github.com/slack-go/slack"
)

type MockClient struct {
	channels map[string]*slack.Channel
}

func (c *MockClient) GetConversationInfo(id string, includeLocale bool) (channel *slack.Channel, err error) {
	return c.channels[id], nil
}

var channels map[string]interface{}

func TestNewConversationFromID(t *testing.T) {
	tc := map[string][]string{
		"Channel":        []string{"CH00001", "#test"},
		"Direct Message": []string{"DM00001", "DM"},
		"Group Message":  []string{"GR00001", "Group messaging with: @some @users"},
	}

	channels := map[string]*slack.Channel{
		"CH00001": &slack.Channel{
			GroupConversation: slack.GroupConversation{
				Name: "test",
				Purpose: slack.Purpose{
					Value: "",
				},
			},
			IsChannel: true,
		},
		"DM00001": &slack.Channel{
			GroupConversation: slack.GroupConversation{},
			IsChannel:         false,
		},
		"GR00001": &slack.Channel{
			GroupConversation: slack.GroupConversation{
				Name: "",
				Purpose: slack.Purpose{
					Value: "Group messaging with: @some @users",
				},
			},
			IsChannel: true,
		},
	}

	client := MockClient{channels: channels}
	for testID, data := range tc {
		conversation, err := NewConversationFromID(data[0], &client)
		if err != nil {
			t.Logf("NewConversationFromID errored for %v: %v", testID, err)
			t.Fail()
		}
		if conversation.Name != data[1] {
			t.Logf("Conversation name was %v, instead of expected %v", conversation.Name, data[1])
			t.Fail()
		}
	}
}
