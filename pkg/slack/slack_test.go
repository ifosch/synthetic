package slack

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	s "github.com/slack-go/slack"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func TestProcess(t *testing.T) {
	disableLogs()
	client := NewMockClient()
	processedMessages := 0
	c := &Chat{
		api:        client,
		processors: map[string][]IMessageProcessor{},
		botID:      "me",
	}
	c.RegisterMessageProcessor(
		NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/slack.CountProcessedMessages",
			func(synthetic.Message) { processedMessages++ },
		),
	)
	messageEvent := s.MessageEvent{
		Msg: s.Msg{
			ClientMsgID:     "CMID001",
			ThreadTimestamp: "",
			User:            "U000001",
			Channel:         "CH00001",
			Text:            "",
		},
		SubMessage:      nil,
		PreviousMessage: nil,
	}
	connectedEvent := s.ConnectedEvent{
		ConnectionCount: 0,
		Info: &s.Info{
			User: &s.UserDetails{
				ID:   "U000002",
				Name: "mybot",
			},
			Team: &s.Team{
				Name: "test-slack",
			},
		},
	}

	c.Process(s.RTMEvent{
		Data: &messageEvent,
	})
	time.Sleep(100 * time.Millisecond)
	c.Process(s.RTMEvent{
		Data: &connectedEvent,
	})

	if processedMessages != 1 {
		t.Logf("Wrong number of processed messages %v should be %v", processedMessages, 1)
		t.Fail()
	}
	if c.botID != "U000002" {
		t.Logf("Wrong botID %v should be U000002", c.botID)
		t.Fail()
	}
}

func TestRegisterMessageProcessor(t *testing.T) {
	disableLogs()
	c := &Chat{
		processors: map[string][]IMessageProcessor{},
	}

	c.RegisterMessageProcessor(
		NewMessageProcessor(
			"github.com/ifosch/synthetic/pkg/slack.LogMessage",
			LogMessage,
		),
	)

	if c.processors["message"][0].Name() != "github.com/ifosch/synthetic/pkg/slack.LogMessage" {
		t.Logf("Wrong processor registered '%v' expected 'github.com/ifosch/synthetic/pkg/slack.LogMessage'", c.processors["message"][0].Name())
		t.Fail()
	}
}

type EventMessageCase struct {
	event    s.MessageEvent
	expected Message
}

func sameConversations(a, b *Conversation) bool {
	if a != nil && b != nil && a.Name() == b.Name() {
		return true
	}
	return false
}

func sameUsers(a, b *User) bool {
	if a != nil && b != nil && a.Name() == b.Name() {
		return true
	}
	return false
}

func sameMessages(a, b *Message) bool {
	if a != nil && b != nil {
		return true
	}
	if a.Completed != b.Completed || a.thread != b.thread || a.mention != b.mention || !sameUsers(a.user, b.user) || !sameConversations(a.conversation, b.conversation) || a.text != b.text {
		return false
	}
	return true
}

func TestReadMessage(t *testing.T) {
	client := NewMockClient()
	chat := &Chat{
		api:                  client,
		rtm:                  nil,
		defaultReplyInThread: false,
		processors:           map[string][]IMessageProcessor{},
		botID:                "me",
	}
	user, _ := NewUserFromID("U000001", client)
	conversation, _ := NewConversationFromID("CH00001", client)
	messageEvents := messageEvents()
	tc := map[string]*EventMessageCase{
		"Incomplete message": {
			event: messageEvents[0],
			expected: &Message{
				Completed:    false,
				thread:       false,
				mention:      false,
				user:         nil,
				conversation: nil,
				text:         "",
			},
		},
		"Threaded message": {
			event: messageEvents[1],
			expected: &Message{
				Completed:    true,
				thread:       true,
				mention:      false,
				user:         user,
				conversation: conversation,
				text:         "",
			},
		},
		"Non-threaded message": {
			event: messageEvents[2],
			expected: &Message{
				Completed:    true,
				thread:       false,
				mention:      false,
				user:         user,
				conversation: conversation,
				text:         "",
			},
		},
		"Message with mention": {
			event: messageEvents[3],
			expected: &Message{
				Completed:    true,
				thread:       false,
				mention:      true,
				user:         user,
				conversation: conversation,
				text:         "<@me>",
			},
		},
	}

	for testID, data := range tc {
		t.Run(testID, func(t *testing.T) {
			message, err := chat.ReadMessage(data.event)
			if err != nil {
				t.Logf("ReadMessage errored for %v: %v", testID, err)
				t.Fail()
			}
			if !sameMessages(message, data.expected) {
				t.Logf("\nMessage in %v test was  %v, \nbut expected %v", testID, message, data.expected)
				t.Fail()
			}
		})
	}
}
