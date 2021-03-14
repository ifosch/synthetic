package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// Start ...
func Start(token string, debug bool) {
	api := slack.New(
		token,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			msg, err := ReadMessage(ev, api)
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
