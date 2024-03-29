package slack

import (
	"fmt"
	"io"
	"log"
	"sync"
	"testing"

	s "github.com/slack-go/slack"
)

func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(io.Discard)
}

func TestProcess(t *testing.T) {
	disableLogs()
	client := NewMockClient()
	processedMessages := 0
	// This is required to avoid race conditions in computing this
	// side effect. Proper solution would be not to rely on side
	// effects for testing.
	//
	// Alternatively we could create a type to hold the counter
	// that controls the locking itself, but that's only useful
	// once we change this test.
	var m sync.Mutex

	c := NewChat(client, false, "me")

	go func(m *sync.Mutex) {
		for {
			_, ok := <-c.MessageChannel
			if !ok {
				t.Error("Borked!")
				continue
			}
			m.Lock()
			processedMessages++
			m.Unlock()
		}
	}(&m)

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
	c.Process(s.RTMEvent{
		Data: &connectedEvent,
	})

	m.Lock()
	if processedMessages != 1 {
		t.Logf("Wrong number of processed messages %v should be %v", processedMessages, 1)
		t.Fail()
	}
	m.Unlock()
	if c.botID != "U000002" {
		t.Logf("Wrong botID %v should be U000002", c.botID)
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

func sameMessages(a, b Message) error {
	switch {
	case a.Completed != b.Completed:
		return fmt.Errorf("value of `Completed` mismatch: %v and %v", a.Completed, b.Completed)
	case a.thread != b.thread:
		return fmt.Errorf("value of `thread` mismatch: %v and %v", a.thread, b.thread)
	case a.mention != b.mention:
		return fmt.Errorf("value of `mention` mismatch: %v and %v", a.mention, b.mention)
	case !sameUsers(a.user, b.user):
		return fmt.Errorf("value of `user` mismatch: %v and %v", a.user, b.user)
	case !sameConversations(a.conversation, b.conversation):
		return fmt.Errorf("value of `conversation` mismatch: %v and %v", a.conversation, b.conversation)
	case a.text != b.text:
		return fmt.Errorf("value of `text` mismatch: `%v` and `%v`", a.text, b.text)
	}
	return nil
}

func TestReadMessage(t *testing.T) {
	client := NewMockClient()
	chat := &Chat{
		api:                  client,
		rtm:                  nil,
		defaultReplyInThread: false,
		botID:                "me",
	}
	user, _ := NewUserFromID("U000001", client)
	conversation, _ := NewConversationFromID("CH00001", client)
	messageEvents := messageEvents()
	tc := map[string]*EventMessageCase{
		// "Incomplete message": {
		// 	event: *messageEvents["uninitialized message"],
		// 	expected: Message{
		// 		Completed:    false,
		// 		thread:       false,
		// 		mention:      false,
		// 		user:         nil,
		// 		conversation: nil,
		// 		text:         "",
		// 	},
		// },
		"Threaded message": {
			event: *messageEvents["empty message in thread"],
			expected: Message{
				Completed:    true,
				thread:       true,
				mention:      false,
				user:         user,
				conversation: conversation,
				text:         "",
			},
		},
		"Non-threaded message": {
			event: *messageEvents["empty message no thread"],
			expected: Message{
				Completed:    true,
				thread:       false,
				mention:      false,
				user:         user,
				conversation: conversation,
				text:         "",
			},
		},
		"Message with mention": {
			event: *messageEvents["only mention no thread"],
			expected: Message{
				Completed:    true,
				thread:       false,
				mention:      true,
				user:         user,
				conversation: conversation,
				text:         "",
			},
		},
	}

	for testID, data := range tc {
		t.Run(testID, func(t *testing.T) {
			message, err := chat.ReadMessage(&data.event)
			if err != nil {
				t.Logf("ReadMessage errored for %v: %v", testID, err)
				t.Fail()
			}
			if err := sameMessages(*message, data.expected); err != nil {
				t.Logf(
					"Message in %v \ntest was: %#v, \nbut expected: %#v",
					testID,
					*message,
					data.expected,
				)
				t.Logf("validation error: %v", err.Error())
				t.Fail()
			}
		})
	}
}
