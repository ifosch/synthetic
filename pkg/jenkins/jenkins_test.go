package jenkins

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	expectedJobs := map[string]string{
		"deploy": "Deploy project",
	}
	mas := NewMockAutomationServer(
		expectedJobs,
	)
	j := &Jenkins{
		jk:   mas,
		jobs: NewJobs(mas),
	}
	j.Load()
	input := "build  deploy      INDEX=users"
	command := "build"

	job, _, _ := j.ParseArgs(input, command)

	if job != "deploy" {
		t.Logf("Wrong job parsed '%v' should be 'deploy'", job)
		t.Fail()
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
	mas := NewMockAutomationServer(
		expectedJobs,
	)
	j := &Jenkins{
		jk:   mas,
		jobs: NewJobs(mas),
	}

	j.Load()

	if j.jobs.Len() != len(expectedJobs) {
		t.Logf("Wrong number of jobs loaded %v but expected %v", j.jobs.Len(), len(expectedJobs))
		t.Fail()
	}
	i := 0
	for job := range j.jobs.jobs {
		if mas.jobs[job] != expectedJobs[job] {
			t.Logf("Wrong job loaded %v expected %v", mas.jobs[job], expectedJobs[job])
			t.Fail()
		}
		i++
	}

	msg := &MockSyntheticMessage{
		replies: []string{},
	}

	j.Reload(msg)

	if j.jobs.Len() != len(expectedJobs) {
		t.Logf("Wrong number of jobs loaded %v but expected %v", j.jobs.Len(), len(expectedJobs))
		t.Fail()
	}
	i = 0
	for job := range j.jobs.jobs {
		if mas.jobs[job] != expectedJobs[job] {
			t.Logf("Wrong job loaded %v expected %v", mas.jobs[job], expectedJobs[job])
			t.Fail()
		}
		i++
	}
	if len(msg.replies) != 1 {
		t.Logf("Wrong number of replies received %v should be 1", len(msg.replies))
		t.Fail()
	}
	expectedReply := "3 Jenkins jobs reloaded"
	if msg.replies[0] != expectedReply {
		t.Logf("Wrong reply '%v' should be '%v'", msg.replies[0], expectedReply)
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
	mas := NewMockAutomationServer(
		expectedJobs,
	)
	j := &Jenkins{
		jk:   mas,
		jobs: NewJobs(mas),
	}
	msg := &MockSyntheticMessage{
		text:    "describe test",
		replies: []string{},
	}
	expectedReply := expectedJobs["test"]

	j.Load()
	j.Describe(msg)

	if len(msg.replies) != 1 {
		t.Logf("Wrong number of replies %v but expected 1", len(msg.replies))
		t.Fail()
	}
	if msg.replies[0] != expectedReply {
		t.Logf("Wrong reply '%v' but expected '%v'", msg.replies[0], expectedReply)
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
	expectedReply := "build\ndeploy\ntest\n"
	mas := NewMockAutomationServer(
		expectedJobs,
	)
	j := &Jenkins{
		jk:   mas,
		jobs: NewJobs(mas),
	}
	msg := &MockSyntheticMessage{
		replies: []string{},
	}

	j.Load()
	j.List(msg)

	if len(msg.replies) != 1 {
		t.Logf("Wrong number of replies %v but expected 1", len(msg.replies))
		t.Fail()
	}
	if msg.replies[0] != expectedReply {
		t.Logf("Wrong reply '%v' but expected '%v'", msg.replies[0], expectedReply)
		t.Fail()
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
	mas := NewMockAutomationServer(
		expectedJobs,
	)
	j := &Jenkins{
		jk:   mas,
		jobs: NewJobs(mas),
	}
	msg := &MockSyntheticMessage{
		text:    "build test",
		replies: []string{},
	}
	expectedReplies := []string{
		"Execution for job `test` was queued",
		fmt.Sprintf("Building `test` with parameters `map[]` (%v/job/test)", os.Getenv("JENKINS_URL")),
		"Job test completed",
	}

	j.Load()
	j.Build(msg)

	if len(msg.replies) != len(expectedReplies) {
		t.Logf("Wrong number of replies %v but expected %v", len(msg.replies), len(expectedReplies))
		t.Fail()
	}
	for i, reply := range msg.replies {
		if reply != expectedReplies[i] {
			t.Logf("Wrong reply '%v' but expected '%v'", reply, expectedReplies[i])
			t.Fail()
		}
	}
}
