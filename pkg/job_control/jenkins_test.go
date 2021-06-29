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
	input         string
	command       string
	expectedJob   string
	expectedArgs  map[string]string
	expectedError string
}

func TestParsing(t *testing.T) {
	disableLogs()
	j := &Jenkins{
		js: NewMockJobServer(
			map[string]string{
				"deploy": "Deploy project",
			},
		),
	}

	tcs := []parsingTC{
		{
			input:         "build  deploy      INDEX=users",
			command:       "build",
			expectedJob:   "deploy",
			expectedArgs:  map[string]string{"INDEX": "users"},
			expectedError: "",
		},
		{
			input:         "build  deploy INDEX=\"users ducks\"",
			command:       "build",
			expectedJob:   "deploy",
			expectedArgs:  map[string]string{"INDEX": "\"users ducks\""},
			expectedError: "",
		},
		{
			input:         "describe",
			command:       "describe",
			expectedJob:   "",
			expectedArgs:  map[string]string{},
			expectedError: "you must specify, at least, one job. You can use `list` to get a list of defined jobs and `describe <job>` to get all details about a job",
		},
		{
			input:         "describe missingjob",
			command:       "describe",
			expectedJob:   "", // Job does not exist so it returns empty
			expectedArgs:  map[string]string{},
			expectedError: "the job `missingjob` doesn't exist in current job list. If it's new addition, try using `reload` to refresh the list of jobs",
		},
	}

	for _, test := range tcs {
		t.Run(test.input, func(t *testing.T) {
			job, args, err := j.ParseArgs(test.input, test.command)

			// Unexpected error happened
			if test.expectedError == "" && err != nil {
				t.Logf("Unexpected error %v", err)
				t.Fail()
			}

			// Expected error did not happen
			if test.expectedError != "" && err == nil {
				t.Logf("Expected error '%v' didn't happen", test.expectedError)
				t.Fail()
			}

			// Job parsing did not match.
			if job != test.expectedJob {
				t.Logf("Wrong job parsed '%v' should be '%v'", job, test.expectedJob)
				t.Fail()
			}

			// Parsed arguments did not match
			for expectedName, expectedValue := range test.expectedArgs {
				value, ok := args[expectedName]
				if !ok {
					t.Logf("Missing argument '%v'", expectedName)
					t.Fail()
				}
				if value != expectedValue {
					t.Logf("Wrong value '%v' for '%v' should be '%v'", value, expectedName, expectedValue)
					t.Fail()
				}
			}
		})
	}
}

type loadTC struct {
	expectedJobs           map[string]string
	expectedReplyOnReload  string
	expectedRepliesOnBuild []string
}

func TestLoadReload(t *testing.T) {
	disableLogs()
	tc := loadTC{
		expectedJobs: map[string]string{
			"build":  "Build the project",
			"test":   "Run test suit on the project",
			"deploy": "Deploy project",
		},
		expectedReplyOnReload: "3 Jenkins jobs reloaded",
	}
	j := &Jenkins{
		js: NewMockJobServer(
			tc.expectedJobs,
		),
	}

	if j.js.GetJobs().Len() != len(tc.expectedJobs) {
		t.Logf("Wrong number of jobs loaded %v but expected %v", j.js.GetJobs().Len(), len(tc.expectedJobs))
		t.Fail()
	}
	i := 0
	for job := range tc.expectedJobs {
		if j.js.GetJob(job).Describe() != tc.expectedJobs[job] {
			t.Logf("Wrong job loaded %v expected %v", j.js.GetJob(job), tc.expectedJobs[job])
			t.Fail()
		}
		i++
	}

	msg := synthetic.NewMockMessage("", false)

	j.Reload(msg)

	if j.js.GetJobs().Len() != len(tc.expectedJobs) {
		t.Logf("Wrong number of jobs loaded %v but expected %v", j.js.GetJobs().Len(), len(tc.expectedJobs))
		t.Fail()
	}
	i = 0
	for job := range tc.expectedJobs {
		if j.js.GetJob(job).Describe() != tc.expectedJobs[job] {
			t.Logf("Wrong job loaded %v expected %v", j.js.GetJob(job).Name(), tc.expectedJobs[job])
			t.Fail()
		}
		i++
	}
	if len(msg.Replies()) != 1 {
		t.Logf("Wrong number of replies received %v should be 1", len(msg.Replies()))
		t.Fail()
	}
	if msg.Replies()[0] != tc.expectedReplyOnReload {
		t.Logf("Wrong reply '%v' should be '%v'", msg.Replies()[0], tc.expectedReplyOnReload)
		t.Fail()
	}
}

func TestDescribe(t *testing.T) {
	disableLogs()
	tc := loadTC{
		expectedJobs: map[string]string{
			"build":  "Build the project",
			"test":   "Run test suit on the project",
			"deploy": "Deploy project",
		},
	}
	j := &Jenkins{
		js: NewMockJobServer(
			tc.expectedJobs,
		),
	}
	msg := synthetic.NewMockMessage("describe test", true)

	j.Describe(msg)

	if len(msg.Replies()) != 1 {
		t.Logf("Wrong number of replies %v but expected 1", len(msg.Replies()))
		t.Fail()
	}
	if msg.Replies()[0] != tc.expectedJobs["test"] {
		t.Logf("Wrong reply '%v' but expected '%v'", msg.Replies()[0], tc.expectedJobs["test"])
		t.Fail()
	}
}

func TestList(t *testing.T) {
	disableLogs()
	tc := loadTC{
		expectedJobs: map[string]string{
			"build":  "Build the project",
			"test":   "Run test suit on the project",
			"deploy": "Deploy project",
		},
	}
	j := &Jenkins{
		js: NewMockJobServer(
			tc.expectedJobs,
		),
	}
	msg := synthetic.NewMockMessage("", false)

	j.List(msg)

	if len(msg.Replies()) != 1 {
		t.Logf("Wrong number of replies %v but expected 1", len(msg.Replies()))
		t.Fail()
	}
	for jobName := range tc.expectedJobs {
		if !strings.Contains(msg.Replies()[0], jobName) {
			t.Logf("Job named '%v' not found in '%v'", jobName, msg.Replies()[0])
			t.Fail()
		}
	}
}

func TestBuild(t *testing.T) {
	disableLogs()
	tc := loadTC{
		expectedJobs: map[string]string{
			"build":  "Build the project",
			"test":   "Run test suit on the project",
			"deploy": "Deploy project",
		},
		expectedRepliesOnBuild: []string{
			"Execution for job `test` was queued",
			fmt.Sprintf("Building `test` with parameters `map[]` (%v/job/test)", os.Getenv("JENKINS_URL")),
			"Job test completed",
		},
	}
	j := &Jenkins{
		js: NewMockJobServer(
			tc.expectedJobs,
		),
	}
	msg := synthetic.NewMockMessage("build test", true)

	j.Build(msg)

	if len(msg.Replies()) != len(tc.expectedRepliesOnBuild) {
		t.Logf("Wrong number of replies %v but expected %v", len(msg.Replies()), len(tc.expectedRepliesOnBuild))
		t.Fail()
	}
	for i, reply := range msg.Replies() {
		if reply != tc.expectedRepliesOnBuild[i] {
			t.Logf("Wrong reply '%v' but expected '%v'", reply, tc.expectedRepliesOnBuild[i])
			t.Fail()
		}
	}
}

func TestTokenizeParams(t *testing.T) {
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
			result := tokenizeParams(tc.input)
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
