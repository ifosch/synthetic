package slack

import (
	"log"
	"os"
	"reflect"
	"runtime"

	"github.com/ifosch/synthetic/pkg/synthetic"
	"github.com/slack-go/slack"
)

// getProcessorName is a helper to calculate the processor name from
// the function name. This is to be changed or removed as soon as
// processors become something better defined.
func getProcessorName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

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
	processors           map[string][]func(synthetic.Message)
	botID                string
}

// NewChat is the constructor for the Chat object.
func NewChat(token string, defaultReplyInThread bool, debug bool) (chat *Chat) {
	api := slack.New(
		token,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	processors := map[string][]func(synthetic.Message){
		"message": []func(synthetic.Message){},
	}
	chat = &Chat{
		api:                  api,
		rtm:                  nil,
		defaultReplyInThread: defaultReplyInThread,
		processors:           processors,
		botID:                "",
	}
	chat.rtm = chat.api.NewRTM()
	chat.RegisterMessageProcessor(LogMessage)
	return
}

// IncomingEvents returns the channel to the chat system events.
func (c *Chat) IncomingEvents() chan slack.RTMEvent {
	return c.rtm.(*slack.RTM).IncomingEvents
}

// RegisterMessageProcessor allows to add more message processors.
func (c *Chat) RegisterMessageProcessor(processor func(synthetic.Message)) {
	c.processors["message"] = append(c.processors["message"], processor)
	log.Printf("%v function registered", getProcessorName(processor))
}

// Start initializes the chat connection.
func (c *Chat) Start() {
	go c.rtm.ManageConnection()

	for msg := range c.IncomingEvents() {
		c.Process(msg)
	}
}

// Process runs the message processing for the chat system.
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
	// Unmanaged events:
	case *slack.AckMessage:
	case *slack.BotAddedEvent:
	case *slack.BotChangedEvent:
	case *slack.ChannelArchiveEvent:
	case *slack.ChannelCreatedEvent:
	case *slack.ChannelDeletedEvent:
	case *slack.ChannelHistoryChangedEvent:
	case *slack.ChannelInfoEvent:
	case *slack.ChannelJoinedEvent:
	case *slack.ChannelLeftEvent:
	case *slack.ChannelMarkedEvent:
	case *slack.ChannelRenameEvent:
	case *slack.ChannelUnarchiveEvent:
	case *slack.CommandsChangedEvent:
	case *slack.DNDUpdatedEvent:
	case *slack.DesktopNotificationEvent:
	case *slack.EmailDomainChangedEvent:
	case *slack.EmojiChangedEvent:
	case *slack.FileChangeEvent:
	case *slack.FileCommentAddedEvent:
	case *slack.FileCommentDeletedEvent:
	case *slack.FileCommentEditedEvent:
	case *slack.FileCreatedEvent:
	case *slack.FileDeletedEvent:
	case *slack.FilePrivateEvent:
	case *slack.FilePublicEvent:
	case *slack.FileSharedEvent:
	case *slack.FileUnsharedEvent:
	case *slack.GroupArchiveEvent:
	case *slack.GroupCloseEvent:
	case *slack.GroupCreatedEvent:
	case *slack.GroupHistoryChangedEvent:
	case *slack.GroupJoinedEvent:
	case *slack.GroupLeftEvent:
	case *slack.GroupMarkedEvent:
	case *slack.GroupOpenEvent:
	case *slack.GroupRenameEvent:
	case *slack.GroupUnarchiveEvent:
	case *slack.HelloEvent:
	case *slack.IMCloseEvent:
	case *slack.IMCreatedEvent:
	case *slack.IMHistoryChangedEvent:
	case *slack.IMMarkedEvent:
	case *slack.IMMarkedHistoryChanged:
	case *slack.IMOpenEvent:
	case *slack.LatencyReport:
	case *slack.ManualPresenceChangeEvent:
	case *slack.MemberJoinedChannelEvent:
	case *slack.MemberLeftChannelEvent:
	case *slack.MobileInAppNotificationEvent:
	case *slack.PinAddedEvent:
	case *slack.PinRemovedEvent:
	case *slack.Ping:
	case *slack.Pong:
	case *slack.PrefChangeEvent:
	case *slack.PresenceChangeEvent:
	case *slack.ReactionAddedEvent:
	case *slack.ReactionRemovedEvent:
	case *slack.ReconnectUrlEvent:
	case *slack.StarAddedEvent:
	case *slack.StarRemovedEvent:
	case *slack.SubteamCreatedEvent:
	case *slack.SubteamMembersChangedEvent:
	case *slack.SubteamSelfAddedEvent:
	case *slack.SubteamSelfRemovedEvent:
	case *slack.SubteamUpdatedEvent:
	case *slack.TeamDomainChangeEvent:
	case *slack.TeamJoinEvent:
	case *slack.TeamMigrationStartedEvent:
	case *slack.TeamPrefChangeEvent:
	case *slack.TeamRenameEvent:
	case *slack.UserChangeEvent:
	case *slack.UserTypingEvent:
	default:
		log.Printf("Unexpected event (%T): %v", ev, ev)
	}
}
