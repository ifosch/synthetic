package synthetic

// Message is an interface for a chat message.
type Message interface {
	Reply(msg string, inThread bool)
	React(reaction string)
	Unreact(reaction string)
	ClearMention() string
	Thread() bool
	Mention() bool
	Text() string
	User() User
	Conversation() Conversation
}
