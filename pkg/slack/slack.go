package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// LogMessage ...
func LogMessage(msg *Message) {
	thread := ""
	if msg.Thread {
		thread = "a thread in "
	}
	log.Printf("Message: '%v' from '%v' in %v'%v'\n", msg.Text, msg.User.Name, thread, msg.Conversation.Name)
}

// Chat is a ...
type Chat struct {
	api                  *slack.Client
	rtm                  *slack.RTM
	defaultReplyInThread bool
	processors           map[string][]func(*Message)
}

// NewChat ...
func NewChat(token string, defaultReplyInThread bool, debug bool) (chat *Chat) {
	api := slack.New(
		token,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	processors := map[string][]func(*Message){
		"message": []func(*Message){},
	}
	chat = &Chat{
		api:                  api,
		rtm:                  nil,
		defaultReplyInThread: defaultReplyInThread,
		processors:           processors,
	}
	chat.rtm = chat.api.NewRTM()
	chat.RegisterMessageProcessor(LogMessage)
	return
}

// RegisterMessageProcessor ...
func (c *Chat) RegisterMessageProcessor(processor func(*Message)) {
	c.processors["message"] = append(c.processors["message"], processor)
	log.Printf("%v function registered", getProcessorName(processor))
}

// Start ...
func (c *Chat) Start() {
	go c.rtm.ManageConnection()

	for msg := range c.rtm.IncomingEvents {
		c.Process(msg)
	}
}

// Process ...
func (c *Chat) Process(msg slack.RTMEvent) {
	switch ev := msg.Data.(type) {
	case *slack.MessageEvent:
		msg, err := ReadMessage(ev, c)
		if err != nil {
			log.Printf("Error %v processing message %v", err, ev)
			break
		}
		if msg.Completed {
			for _, processor := range c.processors["message"] {
				log.Printf("Invoking processor %v", getProcessorName(processor))
				go processor(msg)
			}
		}
	}
}
