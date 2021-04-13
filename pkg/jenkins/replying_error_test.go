package jenkins

import "testing"

func TestReplyingError(t *testing.T) {
	e := ReplyingError{
		Msg: "test",
	}
	if e.Error() != "test" {
		t.Logf("Wrong error message %v expected test", e.Error())
		t.Fail()
	}
}
