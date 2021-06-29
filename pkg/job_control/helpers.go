package jobcontrol

import (
	"strings"
	"unicode"
)

func tokenizeParams(input string) []string {
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
