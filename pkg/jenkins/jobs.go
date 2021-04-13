package jenkins

// Jobs is a collection object to gather all jobs in the Jenkins
// instance.
type Jobs struct {
	jk       AutomationServer
	jobs     map[string]func([]string, map[string]string, chan string)
	jobNames []string
}

// NewJobs is the constructor for the `Jobs` object.
func NewJobs(jk AutomationServer) *Jobs {
	return &Jobs{
		jk:       jk,
		jobs:     map[string]func([]string, map[string]string, chan string){},
		jobNames: []string{},
	}
}

// AddJob adds a new job to the `Jobs` object, providing the name and
// the runner function.
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

// Len provides the amount of jobs.
func (jobs *Jobs) Len() int {
	return len(jobs.jobNames)
}

// String satisfies the `Stringer` interface.
func (jobs *Jobs) String() string {
	jobList := ""
	for _, job := range jobs.jobNames {
		jobList += job + "\n"
	}
	return jobList
}
