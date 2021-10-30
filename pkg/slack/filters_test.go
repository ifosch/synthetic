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
	tcs := map[string]filtersTC{
		"one call": {
			message:       synthetic.NewMockMessage("test", false),
			expectedCalls: 1,
		},
		"no calls": {
			message:       synthetic.NewMockMessage("test ", false),
			expectedCalls: 0,
		},
		"empty message": {
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 0,
		},
	}

	derivedProcessor := Exactly(processor, "test")
	for testID, tc := range tcs {
		t.Run(testID, func(t *testing.T) {
			calls = 0

			derivedProcessor(tc.message)

			if calls != tc.expectedCalls {
				t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
				t.Fail()
			}
		})
	}
}

func TestContains(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := map[string]filtersTC{
		"one call": {
			message:       synthetic.NewMockMessage("test", false),
			expectedCalls: 1,
		},
		"no calls": {
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 0,
		},
	}

	derivedProcessor := Contains(processor, "test")
	for testID, tc := range tcs {
		t.Run(testID, func(t *testing.T) {
			calls = 0

			derivedProcessor(tc.message)

			if calls != tc.expectedCalls {
				t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
				t.Fail()
			}
		})
	}
}

func TestMentioned(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := map[string]filtersTC{
		"no calls": {
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 0,
		},
		"one call": {
			message:       synthetic.NewMockMessage("", true),
			expectedCalls: 1,
		},
	}

	derivedProcessor := Mentioned(processor)
	for testID, tc := range tcs {
		t.Run(testID, func(t *testing.T) {
			calls = 0

			derivedProcessor(tc.message)

			if calls != tc.expectedCalls {
				t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
				t.Fail()
			}
		})
	}
}

func TestNotMentioned(t *testing.T) {
	var calls int
	processor := func(synthetic.Message) {
		calls++
	}
	tcs := map[string]filtersTC{
		"one call": {
			message:       synthetic.NewMockMessage("", false),
			expectedCalls: 1,
		},
		"no calls": {
			message:       synthetic.NewMockMessage("", true),
			expectedCalls: 0,
		},
	}

	derivedProcessor := NotMentioned(processor)
	for testID, tc := range tcs {
		t.Run(testID, func(t *testing.T) {
			calls = 0

			derivedProcessor(tc.message)

			if calls != tc.expectedCalls {
				t.Logf("Wrong number of executions %v should be %v", calls, tc.expectedCalls)
				t.Fail()
			}
		})
	}
}
