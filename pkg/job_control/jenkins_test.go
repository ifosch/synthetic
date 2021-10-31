package jobcontrol

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

type parsingTC struct {
	command       string
	expectedJob   string
	expectedArgs  map[string]string
	expectedError string
}

var expectedJobs = map[string]string{
	"build":  "Build the project",
	"test":   "Run test suit on the project",
	"deploy": "Deploy project",
}

func compareStringLists(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("Wrong number of elements, got %v expected %v", len(a), len(b))
	}
	for i, itemA := range a {
		found := false
		for _, itemB := range b {
			if itemA == itemB {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("Unexpected element '%v', expected '%v'", itemA, b[i])
		}
	}
	return nil
}

func keys(m map[string]string) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func setup() (*Jenkins, synthetic.Message) {
	disableLogs()
	j := NewJenkins("", "", "", NewMockJobServer(expectedJobs))
	msg := synthetic.NewMockMessage("", false)

	return j, msg
}

func TestParsing(t *testing.T) {
	j, _ := setup()

	tcs := map[string]parsingTC{
		"build  deploy      INDEX=users": {
			command:       "build",
			expectedJob:   "deploy",
			expectedArgs:  map[string]string{"INDEX": "users"},
			expectedError: "",
		},
		"build  deploy INDEX=\"users ducks\"": {
			command:       "build",
			expectedJob:   "deploy",
			expectedArgs:  map[string]string{"INDEX": "\"users ducks\""},
			expectedError: "",
		},
		"describe": {
			command:       "describe",
			expectedJob:   "",
			expectedArgs:  map[string]string{},
			expectedError: "you must specify, at least, one job. You can use `list` to get a list of defined jobs and `describe <job>` to get all details about a job",
		},
		"describe missingjob": {
			command:       "describe",
			expectedJob:   "", // Job does not exist so it returns empty
			expectedArgs:  map[string]string{},
			expectedError: "the job `missingjob` doesn't exist in current job list. If it's new addition, try using `reload` to refresh the list of jobs",
		},
	}

	for testID, test := range tcs {
		t.Run(testID, func(t *testing.T) {
			job, args, err := j.ParseArgs(testID, test.command)

			// Unexpected error happened
			if test.expectedError == "" && err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			// Expected error did not happen
			if test.expectedError != "" && err == nil {
				t.Errorf("Expected error '%v' didn't happen", test.expectedError)
			}

			// Job parsing did not match.
			if job != test.expectedJob {
				t.Errorf("Wrong job parsed '%v' should be '%v'", job, test.expectedJob)
			}

			// Parsed arguments did not match
			for expectedName, expectedValue := range test.expectedArgs {
				value, ok := args[expectedName]
				if !ok {
					t.Errorf("Missing argument '%v'", expectedName)
				}
				if value != expectedValue {
					t.Errorf("Wrong value '%v' for '%v' should be '%v'", value, expectedName, expectedValue)
				}
			}
		})
	}
}

func TestLoadReload(t *testing.T) {
	j, _ := setup()

	if err := compareStringLists(keys(j.js.(*MockJobServer).originalJobs), keys(expectedJobs)); err != nil {
		t.Error(err)
	}

	j.js.(*MockJobServer).originalJobs["run"] = "Run an arbitrary command"
	msg := synthetic.NewMockMessage("", false)

	j.Reload(msg)

	if err := compareStringLists(keys(j.js.(*MockJobServer).originalJobs), keys(expectedJobs)); err != nil {
		t.Error(err)
	}
	expectedReply := "4 Jenkins jobs reloaded"
	if err := compareStringLists(msg.Replies(), []string{expectedReply}); err != nil {
		t.Error(err)
	}
}

func TestDescribe(t *testing.T) {
	j, _ := setup()

	for jobName, description := range expectedJobs {
		t.Run(jobName, func(t *testing.T) {
			msg := synthetic.NewMockMessage(fmt.Sprintf("describe %s", jobName), true)

			j.Describe(msg)

			if msg.Replies()[0] != description {
				t.Errorf("Wrong description received '%s', expected '%s'", msg.Replies()[0], description)
			}
		})
	}
}

func TestList(t *testing.T) {
	j, msg := setup()

	j.List(msg)

	if len(msg.(*synthetic.MockMessage).Replies()) != 1 {
		t.Errorf("Wrong number of replies %v but expected 1", len(msg.(*synthetic.MockMessage).Replies()))
	}
	jobNames := strings.Split(msg.(*synthetic.MockMessage).Replies()[0], "\n")
	if err := compareStringLists(jobNames[:len(jobNames)-1], keys(expectedJobs)); err != nil {
		t.Error(err)
	}
}

func TestBuild(t *testing.T) {
	j, msg := setup()
	msg.(*synthetic.MockMessage).OriginalText = "build test"
	expectedRepliesOnBuild := []string{
		"Execution for job `test` was queued",
		fmt.Sprintf("Building `test` with parameters `map[]` (%v/job/test)", os.Getenv("JENKINS_URL")),
		"Job test completed",
	}

	j.Build(msg)

	if err := compareStringLists(msg.(*synthetic.MockMessage).Replies(), expectedRepliesOnBuild); err != nil {
		t.Error(err)
	}
}

func TestTokenizeParams(t *testing.T) {
	tt := map[string]struct {
		result []string
	}{
		"": {
			result: []string{},
		},
		"build deploy": {
			result: []string{"build", "deploy"},
		},
		"build  deploy      INDEX=users": {
			result: []string{"build", "deploy", "INDEX=users"},
		},
		"build  deploy      INDEX=\"users\"": {
			result: []string{"build", "deploy", "INDEX=\"users\""},
		},
		"build  deploy      INDEX=\"users ducks\"": {
			result: []string{"build", "deploy", "INDEX=\"users ducks\""},
		},
	}
	for testID, tc := range tt {
		t.Run(testID, func(t *testing.T) {
			result := tokenizeParams(testID)

			if err := compareStringLists(result, tc.result); err != nil {
				t.Error(err)
			}
		})
	}
}
