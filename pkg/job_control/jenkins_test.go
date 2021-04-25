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

type parsingTC struct {
	input           string
	command         string
	expectedJob     string
	expectedArgs    map[string]string
	expectedReplies []string
}

func TestParsing(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	j := &Jenkins{
		js: NewMockJobServer(
			map[string]string{
				"deploy": "Deploy project",
			},
		),
	}

	tcs := []parsingTC{
		{
			input:           "build  deploy      INDEX=users",
			command:         "build",
			expectedJob:     "deploy",
			expectedArgs:    map[string]string{"INDEX": "users"},
			expectedReplies: []string{},
		},
		{
			input:        "describe",
			command:      "describe",
			expectedJob:  "",
			expectedArgs: map[string]string{},
			expectedReplies: []string{
				"You must specify, at least, one job. You can use `list` to get a list of defined jobs and `describe <job>` to get all details about a job.",
			},
		},
		{
			input:        "describe missingjob",
			command:      "describe",
			expectedJob:  "missingjob",
			expectedArgs: map[string]string{},
			expectedReplies: []string{
				"The job `missingjob` doesn't exist in current job list. If it's new addition, try using `reload` to refresh the list of jobs.",
			},
		},
	}

	for _, test := range tcs {
		msg := synthetic.NewMockMessage(test.input)

		job, args := j.ParseArgs(msg, test.command)

		if job != test.expectedJob {
			t.Logf("Wrong job parsed '%v' should be '%v'", job, test.expectedJob)
			t.Fail()
		}
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
		if len(msg.Replies()) != len(test.expectedReplies) {
			t.Logf("Wrong number of replies '%v' should be '%v'", len(msg.Replies()), len(test.expectedReplies))
			t.Fail()
		}
		for i, reply := range msg.Replies() {
			if reply != test.expectedReplies[i] {
				t.Logf("Wrong reply '%v' should be '%v'", reply, test.expectedReplies[i])
				t.Fail()
			}
		}
	}
}

func TestLoadReload(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	expectedJobs := map[string]string{
		"build":  "Build the project",
		"test":   "Run test suit on the project",
		"deploy": "Deploy project",
	}
	j := &Jenkins{
		js: NewMockJobServer(
			expectedJobs,
		),
	}

	if j.js.GetJobs().Len() != len(expectedJobs) {
		t.Logf("Wrong number of jobs loaded %v but expected %v", j.js.GetJobs().Len(), len(expectedJobs))
		t.Fail()
	}
	i := 0
	for job := range expectedJobs {
		if j.js.GetJob(job).Describe() != expectedJobs[job] {
			t.Logf("Wrong job loaded %v expected %v", j.js.GetJob(job), expectedJobs[job])
			t.Fail()
		}
		i++
	}

	msg := synthetic.NewMockMessage("")

	j.Reload(msg)

	if j.js.GetJobs().Len() != len(expectedJobs) {
		t.Logf("Wrong number of jobs loaded %v but expected %v", j.js.GetJobs().Len(), len(expectedJobs))
		t.Fail()
	}
	i = 0
	for job := range expectedJobs {
		if j.js.GetJob(job).Describe() != expectedJobs[job] {
			t.Logf("Wrong job loaded %v expected %v", j.js.GetJob(job).Name(), expectedJobs[job])
			t.Fail()
		}
		i++
	}
	if len(msg.Replies()) != 1 {
		t.Logf("Wrong number of replies received %v should be 1", len(msg.Replies()))
		t.Fail()
	}
	expectedReply := "3 Jenkins jobs reloaded"
	if msg.Replies()[0] != expectedReply {
		t.Logf("Wrong reply '%v' should be '%v'", msg.Replies()[0], expectedReply)
		t.Fail()
	}
}

func TestDescribe(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	expectedJobs := map[string]string{
		"build":  "Build the project",
		"test":   "Run test suit on the project",
		"deploy": "Deploy project",
	}
	j := &Jenkins{
		js: NewMockJobServer(
			expectedJobs,
		),
	}
	msg := synthetic.NewMockMessage("describe test")
	expectedReply := expectedJobs["test"]

	j.Describe(msg)

	if len(msg.Replies()) != 1 {
		t.Logf("Wrong number of replies %v but expected 1", len(msg.Replies()))
		t.Fail()
	}
	if msg.Replies()[0] != expectedReply {
		t.Logf("Wrong reply '%v' but expected '%v'", msg.Replies()[0], expectedReply)
		t.Fail()
	}
}

func TestList(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	expectedJobs := map[string]string{
		"build":  "Build the project",
		"test":   "Run test suit on the project",
		"deploy": "Deploy project",
	}
	j := &Jenkins{
		js: NewMockJobServer(
			expectedJobs,
		),
	}
	msg := synthetic.NewMockMessage("")

	j.List(msg)

	if len(msg.Replies()) != 1 {
		t.Logf("Wrong number of replies %v but expected 1", len(msg.Replies()))
		t.Fail()
	}
	for jobName := range expectedJobs {
		if !strings.Contains(msg.Replies()[0], jobName) {
			t.Logf("Job named '%v' not found in '%v'", jobName, msg.Replies()[0])
			t.Fail()
		}
	}
}

func TestBuild(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	expectedJobs := map[string]string{
		"build":  "Build the project",
		"test":   "Run test suit on the project",
		"deploy": "Deploy project",
	}
	j := &Jenkins{
		js: NewMockJobServer(
			expectedJobs,
		),
	}
	msg := synthetic.NewMockMessage("build test")
	expectedReplies := []string{
		"Execution for job `test` was queued",
		fmt.Sprintf("Building `test` with parameters `map[]` (%v/job/test)", os.Getenv("JENKINS_URL")),
		"Job test completed",
	}

	j.Build(msg)

	if len(msg.Replies()) != len(expectedReplies) {
		t.Logf("Wrong number of replies %v but expected %v", len(msg.Replies()), len(expectedReplies))
		t.Fail()
	}
	for i, reply := range msg.Replies() {
		if reply != expectedReplies[i] {
			t.Logf("Wrong reply '%v' but expected '%v'", reply, expectedReplies[i])
			t.Fail()
		}
	}
}
