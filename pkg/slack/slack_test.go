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
