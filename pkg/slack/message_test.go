package slack

import (
	"testing"

	"github.com/slack-go/slack"
)

func messageEvents() map[string]*slack.MessageEvent {
	return map[string]*slack.MessageEvent{
		"uninitialized message": {
			Msg: slack.Msg{
				ClientMsgID:     "",
				ThreadTimestamp: "",
				User:            "",
				Channel:         "",
				Text:            "",
			},
		},
		"empty message in thread": {
			Msg: slack.Msg{
				ClientMsgID:     "M000001",
				ThreadTimestamp: "165783949832",
				User:            "U000001",
				Channel:         "CH00001",
				Text:            "",
			},
		},
		"empty message no thread": {
			Msg: slack.Msg{
				ClientMsgID:     "M000001",
				ThreadTimestamp: "",
				User:            "U000001",
				Channel:         "CH00001",
				Text:            "",
			},
		},
		"only mention no thread": {
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

func TestReply(t *testing.T) {
	client := NewMockClient()
	rtm := NewMockRTM()
	chat := &Chat{
		api:                  client,
		rtm:                  rtm,
		defaultReplyInThread: false,
		botID:                "me",
	}
	messageEvents := messageEvents()

	for testID, messageEvent := range messageEvents {
		t.Run(testID, func(t *testing.T) {
			message, err := chat.ReadMessage(messageEvent)
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
		})
	}
}

func TestReactUnreact(t *testing.T) {
	client := NewMockClient()
	chat := NewChat(client, false, "me")
	messageEvents := messageEvents()

	for testID, messageEvent := range messageEvents {
		t.Run(testID, func(t *testing.T) {
			message, err := chat.ReadMessage(messageEvent)
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
		})
	}
}
