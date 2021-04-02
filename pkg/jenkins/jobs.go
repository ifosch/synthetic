package jenkins

import (
	jk "github.com/ifosch/jk/pkg/jenkins"
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

// JobIsPresent checks for presence of a `name` job in the `Jobs`
// list.
func (jobs *Jobs) JobIsPresent(name string) bool {
	found := false
	for _, i := range jobs.jobNames {
		if i == name {
			found = true
			break
		}
	}
	return found
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
