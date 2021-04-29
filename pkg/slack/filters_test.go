package slack

import (
	"testing"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

func TestExactly(t *testing.T) {
	calls := 0
	processor := func(synthetic.Message) {
		calls++
	}

	derivedProcessor := Exactly(processor, "test")
	derivedProcessor(synthetic.NewMockMessage("test", false))
	derivedProcessor(synthetic.NewMockMessage("test ", false))
	derivedProcessor(synthetic.NewMockMessage("", false))

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
	derivedProcessor(synthetic.NewMockMessage("test", false))
	derivedProcessor(synthetic.NewMockMessage("test ", false))
	derivedProcessor(synthetic.NewMockMessage("", false))

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
	derivedProcessor(synthetic.NewMockMessage("test", true))
	derivedProcessor(synthetic.NewMockMessage("", false))

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
	derivedProcessor(synthetic.NewMockMessage("test", true))
	derivedProcessor(synthetic.NewMockMessage("", false))

	if calls != 1 {
		t.Logf("Wrong number of executions %v should be 1", calls)
		t.Fail()
	}
}
