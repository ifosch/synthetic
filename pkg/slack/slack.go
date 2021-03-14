package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// Chat is a ...
type Chat struct {
	api *slack.Client
}

// NewChat ...
func NewChat(token string, debug bool) (chat *Chat) {
	api := slack.New(
		token,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	return &Chat{api: api}
}

// Start ...
func (c *Chat) Start() {
	rtm := c.api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			msg, err := ReadMessage(ev, c.api)
			if err != nil {
				log.Printf("Error %v processing message %v", err, ev)
				break
			}
			if msg.Completed {
				thread := ""
				if msg.Thread {
					thread = "a thread in "
				}
				log.Printf("Message: '%v' from '%v' in %v'%v'\n", msg.Text, msg.User.Name, thread, msg.Conversation.Name)
			}
		}
	}
}
