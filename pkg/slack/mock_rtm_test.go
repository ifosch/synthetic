package slack

import (
	"github.com/slack-go/slack"
)

type MockRTM struct {
	messagesSent []*slack.OutgoingMessage
}

func (rtm *MockRTM) ManageConnection() {
}

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

func (rtm *MockRTM) SendMessage(msg *slack.OutgoingMessage) {
	rtm.messagesSent = append(rtm.messagesSent, msg)
}

func (rtm *MockRTM) reset() {
	rtm.messagesSent = []*slack.OutgoingMessage{}
}

func NewMockRTM() *MockRTM {
	rtm := &MockRTM{}
	rtm.reset()
	return rtm
}
