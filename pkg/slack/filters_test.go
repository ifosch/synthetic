package slack

import (
	"testing"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

type filtersTC struct {
	message       *synthetic.MockMessage
	expectedCalls int
}

func TestExactly(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := []filtersTC{
		{
			message:       synthetic.NewMockMessage("test", false),
			expectedCalls: 1,
		},
		{
			message:       synthetic.NewMockMessage("test ", false),
			expectedCalls: 0,
		},
		{
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 0,
		},
	}

	derivedProcessor := Exactly(processor, "test")
	for _, tc := range tcs {
		calls = 0

		derivedProcessor(tc.message)

		if calls != tc.expectedCalls {
			t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
			t.Fail()
		}
	}
}

func TestContains(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := []filtersTC{
		{
			message:       synthetic.NewMockMessage("test", false),
			expectedCalls: 1,
		},
		{
			message:       synthetic.NewMockMessage("test ", false),
			expectedCalls: 1,
		},
		{
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 0,
		},
	}

	derivedProcessor := Contains(processor, "test")
	for _, tc := range tcs {
		calls = 0

		derivedProcessor(tc.message)

		if calls != tc.expectedCalls {
			t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
			t.Fail()
		}
	}
}

func TestMentioned(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := []filtersTC{
		{
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 0,
		},
		{
			message:       synthetic.NewMockMessage("", true),
			expectedCalls: 1,
		},
	}

	derivedProcessor := Mentioned(processor)
	for _, tc := range tcs {
		calls = 0

		derivedProcessor(tc.message)

		if calls != tc.expectedCalls {
			t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
			t.Fail()
		}
	}
}

func TestNotMentioned(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := []filtersTC{
		{
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 1,
		},
		{
			message:       synthetic.NewMockMessage("", true),
			expectedCalls: 0,
		},
	}

	derivedProcessor := NotMentioned(processor)
	for _, tc := range tcs {
		calls = 0

		derivedProcessor(tc.message)

		if calls != tc.expectedCalls {
			t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
			t.Fail()
		}
	}
}
