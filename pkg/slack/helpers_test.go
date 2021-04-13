package slack

import (
	"testing"
)

func TestRemoveWord(t *testing.T) {
	tc := map[string][]string{
		"Simple removal": {"This is a sentence", "This", "is a sentence"},
	}

	for testID, data := range tc {
		result := RemoveWord(data[0], data[1])
		if result != data[2] {
			t.Logf("%v: Removing '%v' from '%v' returned '%v', but '%v' was expected", testID, data[1], data[0], result, data[2])
			t.Fail()
		}
	}
}
