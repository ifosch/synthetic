package jenkins

import (
	jk "github.com/ifosch/jk/pkg/jenkins"
	"github.com/ifosch/synthetic-david/pkg/slack"
)

// Jobs ...
type Jobs struct {
	jk       *jk.Jenkins
	jobs     map[string]func([]string, map[string]string, chan string)
	jobNames []string
}

// NewJobs ...
func NewJobs(jk *jk.Jenkins) *Jobs {
	return &Jobs{
		jk:       jk,
		jobs:     map[string]func([]string, map[string]string, chan string){},
		jobNames: []string{},
	}
}

// AddJob ...
func (jobs *Jobs) AddJob(name string, runner func([]string, map[string]string, chan string)) {
	jobs.jobs[name] = runner
	jobs.jobNames = append(jobs.jobNames, name)
}

// JobIsPresent ...
func (jobs *Jobs) JobIsPresent(name string) bool {
	return slack.InStringSlice(jobs.jobNames, name)
}

// Len ...
func (jobs *Jobs) Len() int {
	return len(jobs.jobNames)
}

// String ...
func (jobs *Jobs) String() string {
	jobList := ""
	for job := range jobs.jobs {
		jobList += job + "\n"
	}
	return jobList
}
