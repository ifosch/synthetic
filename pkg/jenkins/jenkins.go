package jenkins

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	jk "github.com/ifosch/jk/pkg/jenkins"
	"github.com/ifosch/synthetic/pkg/slack"
	"github.com/ifosch/synthetic/pkg/synthetic"
)

func runJob(jobName string, j AutomationServer) func([]string, map[string]string, chan string) {
	return func(params []string, args map[string]string, responses chan string) {
		log.Println(jobName, "job run: ", params, args)
		out := make(chan jk.Message)
		defer close(out)
		log.Println("Asking", j, "to run", jobName, "with args", args)
		go j.Build(jobName, args, out)
		var resp jk.Message
		for {
			resp = <-out
			responses <- resp.Message
			log.Println(resp.Message)
			if resp.Done {
				break
			}
		}
	}
}

// Jenkins is the object to handle the Jenkins connection.
type Jenkins struct {
	jk   AutomationServer
	jobs *Jobs
}

// Init connects to Jenkins and gathers the data from it.
func (j *Jenkins) Init() (err error) {
	j.jk, err = jk.Connect()
	j.jobs = NewJobs(j.jk)
	return
}

// ParseArgs provides parameters and options parsing from a string.
func (j *Jenkins) ParseArgs(input, command string) (job string, params []string, args map[string]string, err error) {
	params = strings.Split(slack.RemoveWord(input, command), " ")
	args = map[string]string{}
	toRemove := []int{}
	for i, param := range params {
		if strings.Contains(param, "=") {
			data := strings.Split(param, "=")
			args[data[0]] = data[1]
			toRemove = append(toRemove, i)
		}
	}
	params = RemoveIndexes(params, toRemove)

	if params[0] == "" {
		err = ReplyingError{
			Msg: fmt.Sprintf(
				"You must specify, at least, one job to run. You can use `list` to get a list of defined jobs. Some jobs might require arguments to run. You can use `describe <job>` to get all details about a job."),
		}
		return
	}

	if !j.jobs.JobIsPresent(params[0]) {
		err = ReplyingError{
			Msg: fmt.Sprintf(
				"The job `%v` doesn't exist in current job list. If it's new addition, try using `reload` to refresh the list of jobs.", params[0]),
		}
		return
	}

	return params[0], params[1:], args, nil
}

// Load loads all elements from Jenkins.
func (j *Jenkins) Load() {
	log.Println("Loading jobs")
	out := make(chan jk.Message)
	defer close(out)
	go j.jk.List(out)
	var update jk.Message
	for {
		update = <-out
		if update.Message != "" {
			log.Println("Registering", update.Message, "job")
			j.jobs.AddJob(update.Message, runJob(update.Message, j.jk))
		}
		if update.Done {
			break
		}
	}
	log.Printf("Loaded %v jobs\n", j.jobs.Len())
}

// Describe replies `msg` with the description of a job defined.
func (j *Jenkins) Describe(msg synthetic.Message) {
	job, _, _, err := j.ParseArgs(msg.ClearMention(), "describe")
	if err != nil {
		log.Println("Error", err, "parsing", msg.Text())
		if errors.As(err, &ReplyingError{}) {
			msg.Reply(err.Error(), msg.Thread())
		}
		return
	}

	responses := make(chan jk.Message)
	defer close(responses)
	go j.jk.Describe(job, nil, responses)
	var resp jk.Message
	description := ""
	for {
		resp = <-responses
		description += resp.Message
		if resp.Done {
			break
		}
	}
	msg.Reply(description, msg.Thread())
}

// List replies `msg` with the list of jobs in the Jenkins instance.
func (j *Jenkins) List(msg synthetic.Message) {
	msg.Reply(fmt.Sprintf("%v", j.jobs), msg.Thread())
}

// Build runs the runner function for a specified job, with the
// parameters and options specified. It receives the job processing
// updates from Jenkins and reacts and replies with these to `msg`.
func (j *Jenkins) Build(msg synthetic.Message) {
	job, params, args, err := j.ParseArgs(msg.ClearMention(), "build")
	if err != nil {
		log.Println("Error", err, "parsing", msg.Text())
		if errors.As(err, &ReplyingError{}) {
			msg.Reply(err.Error(), msg.Thread())
		}
		return
	}

	msg.React("+1")

	responses := make(chan string)
	defer close(responses)
	go j.jobs.jobs[job](params, args, responses)
	for {
		resp := <-responses
		if strings.Contains(resp, "Build queued") {
			msg.Reply(fmt.Sprintf("Execution for job `%v` was queued", job), msg.Thread())
			msg.React("stopwatch")
		} else if strings.Contains(resp, "Build started") {
			msg.Reply(fmt.Sprintf("Building `%v` with parameters `%v` (%v)", job, args, fmt.Sprintf("%v/job/%v", os.Getenv("JENKINS_URL"), job)), msg.Thread())
			msg.Unreact("stopwatch")
			msg.React("gear")
		} else if strings.Contains(resp, "Build finished") {
			msg.Reply(fmt.Sprintf("Job %v completed", job), msg.Thread())
			msg.Unreact("gear")
			msg.React("heavy_check_mark")
			break
		}
	}
}

// Register loads and initializes the Jenkins object and registers the
// corresponding message processors with `chat`.
func Register(client *slack.Chat) {
	j := Jenkins{}
	err := j.Init()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	j.Load()

	client.RegisterMessageProcessor(slack.Mentioned(slack.Contains(j.List, "list")))
	client.RegisterMessageProcessor(slack.Mentioned(slack.Contains(j.Describe, "describe")))
	client.RegisterMessageProcessor(slack.Mentioned(slack.Contains(j.Build, "build")))
}
