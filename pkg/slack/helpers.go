package slack

import (
	"strings"
)

// RemoveWord removes `word` from `text` and returns the result. Note
// this uses a single space to split the words in `text`.
func RemoveWord(text string, word string) string {
	slice := strings.Split(text, " ")
	i := -1
	for k, v := range slice {
		if v == word {
			i = k
			break
		}
	}
	if i >= 0 {
		slice = append(slice[:i], slice[i+1:]...)
	}
	return strings.Join(slice, " ")
}
