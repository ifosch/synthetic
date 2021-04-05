package slack

import (
	"testing"

	"github.com/ifosch/synthetic/pkg/synthetic"
	"github.com/slack-go/slack"
)

type EventMessageCase struct {
	event    *slack.MessageEvent
	expected *Message
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

func messageEvents() []*slack.MessageEvent {
	return []*slack.MessageEvent{
		&slack.MessageEvent{
			Msg: slack.Msg{
				ClientMsgID:     "",
				ThreadTimestamp: "",
				User:            "",
				Channel:         "",
				Text:            "",
			},
		},
		&slack.MessageEvent{
			Msg: slack.Msg{
				ClientMsgID:     "M000001",
				ThreadTimestamp: "165783949832",
				User:            "U000001",
				Channel:         "CH00001",
				Text:            "",
			},
		},
		&slack.MessageEvent{
			Msg: slack.Msg{
				ClientMsgID:     "M000001",
				ThreadTimestamp: "",
				User:            "U000001",
				Channel:         "CH00001",
				Text:            "",
			},
		},
		&slack.MessageEvent{
			Msg: slack.Msg{
				ClientMsgID:     "M000001",
				ThreadTimestamp: "",
				User:            "U000001",
				Channel:         "CH00001",
				Text:            "<@me>",
			},
		},
	}
}

func TestReadMessage(t *testing.T) {
	client := NewMockClient()
	chat := &Chat{
		api:                  client,
		rtm:                  nil,
		defaultReplyInThread: false,
		processors:           map[string][]func(synthetic.Message){},
		botID:                "me",
	}
	user, _ := NewUserFromID("U000001", client)
	conversation, _ := NewConversationFromID("CH00001", client)
	messageEvents := messageEvents()
	tc := map[string]*EventMessageCase{
		"Incomplete message": &EventMessageCase{
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
		"Threaded message": &EventMessageCase{
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
		"Non-threaded message": &EventMessageCase{
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
		"Message with mention": &EventMessageCase{
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
		message, err := ReadMessage(data.event, chat)
		if err != nil {
			t.Logf("ReadMessage errored for %v: %v", testID, err)
			t.Fail()
		}
		if !sameMessages(message, data.expected) {
			t.Logf("\nMessage in %v test was  %v, \nbut expected %v", testID, message, data.expected)
			t.Fail()
		}
	}
}

func TestReply(t *testing.T) {
	client := NewMockClient()
	rtm := NewMockRTM()
	chat := &Chat{
		api:                  client,
		rtm:                  rtm,
		defaultReplyInThread: false,
		processors:           map[string][]func(synthetic.Message){},
		botID:                "me",
	}
	messageEvents := messageEvents()

	for _, messageEvent := range messageEvents {
		message, err := ReadMessage(messageEvent, chat)
		if err != nil {
			t.Logf("ReadMessage errored: %v", err)
			t.Fail()
		}
		message.Reply("reply", false)
		if len(rtm.messagesSent) != 1 {
			t.Logf("I've sent only one message, but %v were detected", len(rtm.messagesSent))
			t.Fail()
		}
		reply := rtm.messagesSent[0]
		if message.Completed {
			if reply.Channel != message.conversation.slackChannel.ID {
				t.Logf("Wrong channel ID used in reply %v should be %v", reply.Channel, message.conversation.slackChannel.ID)
				t.Fail()
			}
			if reply.Text != "reply" {
				t.Logf("Wrong text in reply %v should be reply", reply.Text)
				t.Fail()
			}
			if reply.ThreadTimestamp != message.event.ThreadTimestamp {
				t.Logf("Wrong timestamp in reply %v should be %v", reply.ThreadTimestamp, message.event.ThreadTimestamp)
				t.Fail()
			}
		} else {
			if reply.Channel != "" {
				t.Logf("Incomplete message should have a nil reply but got %v", reply)
				t.Fail()
			}
		}
		rtm.reset()
	}
}

func TestReactUnreact(t *testing.T) {
	client := NewMockClient()
	rtm := NewMockRTM()
	chat := &Chat{
		api:                  client,
		rtm:                  rtm,
		defaultReplyInThread: false,
		processors:           map[string][]func(synthetic.Message){},
		botID:                "me",
	}
	messageEvents := messageEvents()

	for _, messageEvent := range messageEvents {
		message, err := ReadMessage(messageEvent, chat)
		if err != nil {
			t.Logf("ReadMessage errored: %v", err)
			t.Fail()
		}
		message.React("+1")
		if message.Completed {
			if client.reactionsAdded[0].reaction != "+1" {
				t.Logf("Wrong reaction %v when adding reaction, should be +1", client.reactionsAdded[0].reaction)
				t.Fail()
			}
			if client.reactionsAdded[0].item.Channel != message.event.Channel {
				t.Logf("Wrong channel %v when adding reaction, should be %v", client.reactionsAdded[0].item.Channel, message.event.Channel)
				t.Fail()
			}
			if client.reactionsAdded[0].item.Timestamp != message.event.Timestamp {
				t.Logf("Wrong timestamp %v when adding reaction, should be %v", client.reactionsAdded[0].item.Timestamp, message.event.Timestamp)
				t.Fail()
			}
		}
		message.Unreact("+1")
		if message.Completed {
			if client.reactionsRemoved[0].reaction != "+1" {
				t.Logf("Wrong reaction %v when removing reaction, should be +1", client.reactionsRemoved[0].reaction)
				t.Fail()
			}
			if client.reactionsRemoved[0].item.Channel != message.event.Channel {
				t.Logf("Wrong channel %v when removing reaction, should be %v", client.reactionsRemoved[0].item.Channel, message.event.Channel)
				t.Fail()
			}
			if client.reactionsRemoved[0].item.Timestamp != message.event.Timestamp {
				t.Logf("Wrong timestamp %v when removing reaction, should be %v", client.reactionsRemoved[0].item.Timestamp, message.event.Timestamp)
				t.Fail()
			}
		}
		client.reset()
	}
}
