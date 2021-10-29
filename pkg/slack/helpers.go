package slack

import (
	"strings"
)

// ReplaceSpace ...
func ReplaceSpace(s string) string {
	var result []rune
	const badSpace = '\u00A0'
	for _, r := range s {
		if r == badSpace {
			result = append(result, '\u0020')
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func CleanText(s string) string {
	s = ReplaceSpace(s)
	s = strings.TrimSpace(s)
	return s
}

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
