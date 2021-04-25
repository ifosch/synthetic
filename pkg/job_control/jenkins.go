package jobcontrol

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ifosch/synthetic/pkg/slack"
	"github.com/ifosch/synthetic/pkg/synthetic"
)

// Jenkins is the object to handle the Jenkins connection.
type Jenkins struct {
	js IJobServer
}

// Init connects to Jenkins and gathers the data from it.
func (j *Jenkins) Init() (err error) {
	j.js = &JenkinsJobServer{}
	j.js.Connect(os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASSWORD"))
	return
}

// ParseArgs provides parameters and options parsing from a message.
func (j *Jenkins) ParseArgs(input, command string) (job string, args map[string]string, err error) {
	reduceDupSpaces := regexp.MustCompile(`[ ]{2,}`)
	params := strings.Split(slack.RemoveWord(reduceDupSpaces.ReplaceAllString(input, " "), command), " ")
	args = map[string]string{}
	toRemove := []int{}
	for i, param := range params {
		if strings.Contains(param, "=") {
			data := strings.Split(param, "=")
			args[data[0]] = data[1]
			toRemove = append(toRemove, i)
		}
	}
	job = RemoveIndexes(params, toRemove)[0]

	if job == "" {
		return "", nil, fmt.Errorf("you must specify, at least, one job. You can use `list` to get a list of defined jobs and `describe <job>` to get all details about a job")
	}

	if j.js.GetJob(job) == nil {
		return "", nil, fmt.Errorf("the job `%v` doesn't exist in current job list. If it's new addition, try using `reload` to refresh the list of jobs", job)
	}

	return job, args, nil
}

// Reload runs Load again.
func (j *Jenkins) Reload(msg synthetic.Message) {
	msg.React("+1")
	j.js.GetJobs().Clear()
	err := j.js.Load()
	if err != nil {
		msg.Reply(fmt.Sprintf("Error happened reloading jobs %s", err), msg.Thread())
		return
	}

	msg.Reply(fmt.Sprintf("%v Jenkins jobs reloaded", j.js.GetJobs().Len()), msg.Thread())
	msg.React("heavy_check_mark")
}

// Describe replies `msg` with the description of a job defined.
func (j *Jenkins) Describe(msg synthetic.Message) {
	job, _, err := j.ParseArgs(msg.ClearMention(), "describe")
	if err != nil {
		msg.Reply(fmt.Sprintf("%s", err), msg.Thread())
		return
	}
	msg.Reply(j.js.GetJob(job).Describe(), msg.Thread())
}

// List replies `msg` with the list of jobs in the Jenkins instance.
func (j *Jenkins) List(msg synthetic.Message) {
	msg.Reply(j.js.GetJobs().List(), msg.Thread())
}

// Build runs specified job, with the specified options. It receives
// the job processing updates from Jenkins and reacts and replies with
// these to `msg`.
func (j *Jenkins) Build(msg synthetic.Message) {
	job, args, err := j.ParseArgs(msg.ClearMention(), "build")
	if err != nil {
		msg.Reply(fmt.Sprintf("%s", err), msg.Thread())
		return
	}

	msg.React("+1")

	updates := make(chan Update)
	defer close(updates)

	go j.js.GetJob(job).Run(args, updates)

	lastReaction := ""
	for {
		update := <-updates
		msg.Unreact(lastReaction)
		msg.React(update.Reaction)
		msg.Reply(update.Msg, msg.Thread())
		lastReaction = update.Reaction
		if update.Done {
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

	client.RegisterMessageProcessor(slack.NewMessageProcessor("github.com/ifosch/synthetic/pkg/jenkins.List", slack.Mentioned(slack.Contains(j.List, "list"))))
	client.RegisterMessageProcessor(slack.NewMessageProcessor("github.com/ifosch/synthetic/pkg/jenkins.Describe", slack.Mentioned(slack.Contains(j.Describe, "describe"))))
	client.RegisterMessageProcessor(slack.NewMessageProcessor("github.com/ifosch/synthetic/pkg/jenkins.Build", slack.Mentioned(slack.Contains(j.Build, "build"))))
	client.RegisterMessageProcessor(slack.NewMessageProcessor("github.com/ifosch/synthetic/pkg/jenkins.Reload", slack.Mentioned(slack.Contains(j.Reload, "reload"))))
}
