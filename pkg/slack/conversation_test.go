package slack

import (
	"testing"
)

func TestNewConversationFromID(t *testing.T) {
	tc := map[string][]string{
		"Channel":        {"CH00001", "#test"},
		"Direct Message": {"DM00001", "DM"},
		"Group Message":  {"GR00001", "Group messaging with: @some @users"},
	}

	client := NewMockClient()
	for testID, data := range tc {
		conversation, err := NewConversationFromID(data[0], client)
		if err != nil {
			t.Logf("NewConversationFromID errored for %v: %v", testID, err)
			t.Fail()
		}
		if conversation.name != data[1] {
			t.Logf("Conversation name was %v, instead of expected %v", conversation.name, data[1])
			t.Fail()
		}
	}
}
