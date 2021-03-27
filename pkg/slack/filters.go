package slack

import "strings"

// Contains ...
func Contains(processor func(*Message), catch string) func(*Message) {
	return func(msg *Message) {
		if strings.Contains(msg.Text, catch) {
			processor(msg)
		}
	}
}

// Mentioned ...
func Mentioned(processor func(*Message)) func(*Message) {
	return func(msg *Message) {
		if msg.Mention {
			processor(msg)
		}
	}
}

// NotMentioned ...
func NotMentioned(processor func(*Message)) func(*Message) {
	return func(msg *Message) {
		if !msg.Mention {
			processor(msg)
		}
	}
}
