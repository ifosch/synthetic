package slack

import "strings"

// Contains returns a processor that runs the `processor` if the
// message contains the `catch` string.
func Contains(processor func(*Message), catch string) func(*Message) {
	return func(msg *Message) {
		if strings.Contains(msg.Text, catch) {
			processor(msg)
		}
	}
}

// Mentioned returns a processor that runs the `processor` if the
// message is mentioning the bot.
func Mentioned(processor func(*Message)) func(*Message) {
	return func(msg *Message) {
		if msg.Mention {
			processor(msg)
		}
	}
}

// NotMentioned returns a processor that runs the `processor` if the
// message is not mentioning the bot.
func NotMentioned(processor func(*Message)) func(*Message) {
	return func(msg *Message) {
		if !msg.Mention {
			processor(msg)
		}
	}
}
