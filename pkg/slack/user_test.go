package slack

import (
	"testing"
)

func TestNewUserFromID(t *testing.T) {
	tc := map[string][]string{
		"U000001": []string{"U000001", "@username"},
	}

	client := NewMockClient()
	for testID, data := range tc {
		user, err := NewUserFromID(data[0], client)
		if err != nil {
			t.Logf("NewUserFromID errored for %v: %v", testID, err)
			t.Fail()
		}
		if user.Name != data[1] {
			t.Logf("User name was %v, instead of expected %v", user.Name, data[1])
			t.Fail()
		}
	}
}
