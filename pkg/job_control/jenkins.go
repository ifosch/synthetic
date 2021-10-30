package jobcontrol

import (
	"fmt"
	"strings"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

// Jenkins is the object to handle the Jenkins connection.
type Jenkins struct {
	url, user, password string
	js                  IJobServer
}

// NewJenkins returns a pointer to an initialized Jenkins instance
func NewJenkins(url, user, password string, jobServer IJobServer) *Jenkins {
	j := &Jenkins{url: url, user: user, password: password}
	j.js = jobServer
	return j
}

// Connect to the jenkins server with the credentials used during
// initialization
func (j *Jenkins) Connect() error {
	j.js.Connect(j.url, j.user, j.password)
	return nil
}

// ParseArgs provides parameters and options parsing from a message.
func (j *Jenkins) ParseArgs(input, command string) (job string, args map[string]string, err error) {
	args = make(map[string]string)

	var options []string
	tokens := tokenizeParams(input)
	for _, token := range tokens {
		if token != command {
			if strings.Contains(token, "=") {
				data := strings.Split(token, "=")
				args[data[0]] = data[1]
			} else {
				options = append(options, token)
			}
		}
	}

	if len(options) == 0 {
		return "", nil, fmt.Errorf("you must specify, at least, one job. You can use `list` to get a list of defined jobs and `describe <job>` to get all details about a job")
	}

	job = options[0]

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
	job, _, err := j.ParseArgs(msg.Text(), "describe")
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
	job, args, err := j.ParseArgs(msg.Text(), "build")
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
