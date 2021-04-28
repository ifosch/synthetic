package slack

import (
	"testing"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

type MockMessage struct {
	text    string
	mention bool
}

func (m MockMessage) Reply(msg string, inThread bool) {}
func (m MockMessage) React(reaction string)           {}
func (m MockMessage) Unreact(reaction string)         {}
func (m MockMessage) ClearMention() string {
	return m.Text()
}
func (m MockMessage) Thread() bool {
	return false
}
func (m MockMessage) Mention() bool {
	return m.mention
}
func (m MockMessage) Text() string {
	return m.text
}
func (m MockMessage) User() synthetic.User {
	return nil
}
func (m MockMessage) Conversation() synthetic.Conversation {
	return nil
}

func TestExactly(t *testing.T) {
	calls := 0
	processor := func(synthetic.Message) {
		calls++
	}

	derivedProcessor := Exactly(processor, "test")
	derivedProcessor(MockMessage{text: "test"})
	derivedProcessor(MockMessage{text: "test "})
	derivedProcessor(MockMessage{text: ""})

	if calls != 1 {
		t.Logf("Wrong number of executions %v should be 1", calls)
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	calls := 0
	processor := func(synthetic.Message) {
		calls++
	}

	derivedProcessor := Contains(processor, "test")
	derivedProcessor(MockMessage{text: "test"})
	derivedProcessor(MockMessage{text: "test "})
	derivedProcessor(MockMessage{text: ""})

	if calls != 2 {
		t.Logf("Wrong number of executions %v should be 2", calls)
		t.Fail()
	}
}

func TestMentioned(t *testing.T) {
	calls := 0
	processor := func(synthetic.Message) {
		calls++
	}

	derivedProcessor := Mentioned(processor)
	derivedProcessor(MockMessage{text: "test", mention: true})
	derivedProcessor(MockMessage{text: ""})

	if calls != 1 {
		t.Logf("Wrong number of executions %v should be 1", calls)
		t.Fail()
	}
}

func TestNotMentioned(t *testing.T) {
	calls := 0
	processor := func(synthetic.Message) {
		calls++
	}

	derivedProcessor := NotMentioned(processor)
	derivedProcessor(MockMessage{text: "test", mention: true})
	derivedProcessor(MockMessage{text: ""})

	if calls != 1 {
		t.Logf("Wrong number of executions %v should be 1", calls)
		t.Fail()
	}
}
