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

func TestParsing(t *testing.T) {
	disableLogs()
	j := NewJenkins("", "", "", NewMockJobServer(expectedJobs))

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

type loadTC struct {
	expectedReplyOnReload  string
	expectedRepliesOnBuild []string
}

func TestLoadReload(t *testing.T) {
	disableLogs()
	tc := loadTC{
		expectedReplyOnReload: "3 Jenkins jobs reloaded",
	}
	j := NewJenkins("", "", "", NewMockJobServer(expectedJobs))

	if j.js.GetJobs().Len() != len(expectedJobs) {
		t.Errorf("Wrong number of jobs loaded %v but expected %v", j.js.GetJobs().Len(), len(expectedJobs))
	}
	i := 0
	for job := range expectedJobs {
		if j.js.GetJob(job).Describe() != expectedJobs[job] {
			t.Errorf("Wrong job loaded %v expected %v", j.js.GetJob(job), expectedJobs[job])
		}
		i++
	}

	msg := synthetic.NewMockMessage("", false)

	j.Reload(msg)

	if j.js.GetJobs().Len() != len(expectedJobs) {
		t.Errorf("Wrong number of jobs loaded %v but expected %v", j.js.GetJobs().Len(), len(expectedJobs))
	}
	i = 0
	for job := range expectedJobs {
		if j.js.GetJob(job).Describe() != expectedJobs[job] {
			t.Errorf("Wrong job loaded %v expected %v", j.js.GetJob(job).Name(), expectedJobs[job])
		}
		i++
	}
	if len(msg.Replies()) != 1 {
		t.Errorf("Wrong number of replies received %v should be 1", len(msg.Replies()))
	}
	if msg.Replies()[0] != tc.expectedReplyOnReload {
		t.Errorf("Wrong reply '%v' should be '%v'", msg.Replies()[0], tc.expectedReplyOnReload)
	}
}

func TestDescribe(t *testing.T) {
	disableLogs()
	j := NewJenkins("", "", "", NewMockJobServer(expectedJobs))
	msg := synthetic.NewMockMessage("describe test", true)

	j.Describe(msg)

	if len(msg.Replies()) != 1 {
		t.Errorf("Wrong number of replies %v but expected 1", len(msg.Replies()))
	}
	if msg.Replies()[0] != expectedJobs["test"] {
		t.Errorf("Wrong reply '%v' but expected '%v'", msg.Replies()[0], expectedJobs["test"])
	}
}

func TestList(t *testing.T) {
	disableLogs()
	j := NewJenkins("", "", "", NewMockJobServer(expectedJobs))
	msg := synthetic.NewMockMessage("", false)

	j.List(msg)

	if len(msg.Replies()) != 1 {
		t.Errorf("Wrong number of replies %v but expected 1", len(msg.Replies()))
	}
	for jobName := range expectedJobs {
		if !strings.Contains(msg.Replies()[0], jobName) {
			t.Errorf("Job named '%v' not found in '%v'", jobName, msg.Replies()[0])
		}
	}
}

func TestBuild(t *testing.T) {
	disableLogs()
	tc := loadTC{
		expectedRepliesOnBuild: []string{
			"Execution for job `test` was queued",
			fmt.Sprintf("Building `test` with parameters `map[]` (%v/job/test)", os.Getenv("JENKINS_URL")),
			"Job test completed",
		},
	}
	j := NewJenkins("", "", "", NewMockJobServer(expectedJobs))
	msg := synthetic.NewMockMessage("build test", true)

	j.Build(msg)

	if len(msg.Replies()) != len(tc.expectedRepliesOnBuild) {
		t.Errorf("Wrong number of replies %v but expected %v", len(msg.Replies()), len(tc.expectedRepliesOnBuild))
	}
	for i, reply := range msg.Replies() {
		if reply != tc.expectedRepliesOnBuild[i] {
			t.Errorf("Wrong reply '%v' but expected '%v'", reply, tc.expectedRepliesOnBuild[i])
		}
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
