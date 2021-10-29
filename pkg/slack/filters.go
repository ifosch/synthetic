package slack

import (
	"strings"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

// Exactly returns a processot that runs the `processor` if the
// message is exactly like the `catch` string.
func Exactly(processor func(synthetic.Message), catch string) func(synthetic.Message) {
	return func(msg synthetic.Message) {
		if msg.Text() == catch {
			processor(msg)
		}
	}
}

// Contains returns a processor that runs the `processor` if the
// message contains the `catch` string.
func Contains(processor func(synthetic.Message), catch string) func(synthetic.Message) {
	return func(msg synthetic.Message) {
		if strings.Contains(msg.Text(), catch) {
			processor(msg)
		}
	}
}

// Mentioned returns a processor that runs the `processor` if the
// message is mentioning the bot.
func Mentioned(processor func(synthetic.Message)) func(synthetic.Message) {
	return func(msg synthetic.Message) {
		if msg.Mention() {
			processor(msg)
		}
	}
}

// NotMentioned returns a processor that runs the `processor` if the
// message is not mentioning the bot.
func NotMentioned(processor func(synthetic.Message)) func(synthetic.Message) {
	return func(msg synthetic.Message) {
		if !msg.Mention() {
			processor(msg)
		}
	}
}
