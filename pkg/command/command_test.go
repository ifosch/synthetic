package command

import "testing"

func TestTokenizeCommand(t *testing.T) {
	tt := []struct {
		input  string
		result []string
	}{
		{
			input:  "",
			result: []string{},
		},
		{
			input:  "build deploy",
			result: []string{"build", "deploy"},
		},
		{
			input:  "build  deploy      INDEX=users",
			result: []string{"build", "deploy", "INDEX=users"},
		},
		{
			input:  "build  deploy      INDEX=\"users\"",
			result: []string{"build", "deploy", "INDEX=\"users\""},
		},
		{
			input:  "build  deploy      INDEX=\"users ducks\"",
			result: []string{"build", "deploy", "INDEX=\"users ducks\""},
		},
	}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			result := tokenizeCommand(tc.input)
			if len(result) != len(tc.result) {
				t.Errorf("expected %d results but got %d", len(tc.result), len(result))
			}

			for i, value := range result {
				if value != tc.result[i] {
					t.Errorf("expected element %d to be %s but was %s", i, tc.result[i], value)
				}
			}
		})
	}
}
