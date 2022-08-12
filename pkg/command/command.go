package command

import (
	"strings"
	"unicode"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

// Command represents an instance of a command from a message, ready
// to be executed
type Command struct {
	tokenizedParams []string
	message         synthetic.Message
}

// NewCommand creates a new instance of Command based on a message
func NewCommand(message synthetic.Message) *Command {
	return &Command{
		message:         message,
		tokenizedParams: tokenizeCommand(message.Text()),
	}
}

// Message returns the message from the command
func (c *Command) Message() synthetic.Message {
	return c.message
}

// Parses the message and returns a list of tokens for the command
// logic to act on it.
func tokenizeCommand(input string) []string {
	// Solution taken from
	// https://groups.google.com/g/golang-nuts/c/pNwqLyfl2co/m/APaZSSvQUAAJ
	lastQuote := rune(0)
	// Must return true for symbols that delimit a field
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			// when the quotation ends
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			//when we are inside a quotation
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			// when c is a valid quotation mark symbol starts a quotation
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}
	return strings.FieldsFunc(input, f)
}
