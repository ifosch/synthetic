package slack

import (
	"log"
	"os"

	"github.com/ifosch/synthetic/pkg/synthetic"
	"github.com/slack-go/slack"
)

// LogMessage is a message processor to log the message received.
func LogMessage(msg synthetic.Message) {
	thread := ""
	if msg.Thread() {
		thread = "a thread in "
	}
	log.Printf("Message: '%v' from '%v' in %v'%v'\n", msg.Text(), msg.User().Name(), thread, msg.Conversation().Name())
}

// Chat represents the whole chat connection providing methods to
// interact with the chat system.
type Chat struct {
	api                  IClient
	rtm                  IRTM
	defaultReplyInThread bool
	processors           map[string][]IMessageProcessor
	botID                string
}

// NewChat is the constructor for the Chat object.
func NewChat(token string, defaultReplyInThread bool, debug bool) (chat *Chat) {
	api := slack.New(
		token,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	processors := map[string][]IMessageProcessor{
		"message": {},
	}
	chat = &Chat{
		api:                  api,
		rtm:                  nil,
		defaultReplyInThread: defaultReplyInThread,
		processors:           processors,
		botID:                "",
	}
	chat.rtm = chat.api.NewRTM()
	chat.RegisterMessageProcessor(
		NewMessageProcessor("github.com/ifosch/pkg/slack.LogMessage", LogMessage),
	)
	return
}

// IncomingEvents returns the channel to the chat system events.
func (c *Chat) IncomingEvents() chan slack.RTMEvent {
	return c.rtm.(*slack.RTM).IncomingEvents
}

// RegisterMessageProcessor allows to add more message processors.
func (c *Chat) RegisterMessageProcessor(processor IMessageProcessor) {
	c.processors["message"] = append(c.processors["message"], processor)
	log.Printf("%v function registered", processor.Name())
}

// Start initializes the chat connection.
func (c *Chat) Start() {
	go c.rtm.ManageConnection()

	for msg := range c.IncomingEvents() {
		c.Process(msg)
	}
}

func processMessage(ev *slack.MessageEvent, c *Chat) {
	msg, err := ReadMessage(ev, c)
	if err != nil {
		log.Printf("Error %v processing message %v", err, ev)
		return
	}
	if msg.Completed {
		for _, processor := range c.processors["message"] {
			log.Printf("Invoking processor %v", processor.Name())
			go processor.Run(msg)
		}
	}
}

// Process runs the message processing for the chat system.
func (c *Chat) Process(msg slack.RTMEvent) {
	switch ev := msg.Data.(type) {
	case *slack.MessageEvent:
		processMessage(ev, c)
	case *slack.ConnectingEvent:
		log.Printf("Trying to connect to Slack: Attempt %v of %v", ev.Attempt, ev.ConnectionCount)
	case *slack.ConnectedEvent:
		c.botID = ev.Info.User.ID
		log.Printf("Connected to %v Slack as %v after %v attempts ", ev.Info.Team.Name, ev.Info.User.Name, ev.ConnectionCount+1)
	case *slack.InvalidAuthEvent:
		log.Fatalf("Invalid credentials provided to Slack")
	case *slack.ConnectionErrorEvent:
		log.Printf("Error connecting to Slack %v", ev)
	case *slack.DisconnectedEvent:
		log.Printf("Disconnected event: %v", ev)
	case *slack.IncomingEventError:
		log.Printf("Unexpected error receiving a websocket event: %v", ev)
	case *slack.MessageTooLongEvent:
		log.Printf("Last message was too long: %v", ev)
	case *slack.OutgoingErrorEvent:
		log.Printf("Unspecific error on outgoing message: %v", ev)
	case *slack.RTMError:
		log.Printf("Unspecific error on RTM: %v", ev)
	case *slack.RateLimitEvent:
		log.Printf("Slack rate limit reached: %v", ev)
	case *slack.UnmarshallingErrorEvent:
		log.Printf("Unmarshalling error: %v", ev)
	default:
		log.Printf("Unmanaged event (%T)", ev)
	}
}
