package slack

import (
	"reflect"
	"runtime"
	"strings"
)

// getProcessorName ...
func getProcessorName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// RemoveWord ...
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
