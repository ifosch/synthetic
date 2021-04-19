package jobcontrol

import (
	"github.com/bndr/gojenkins"
)

// JenkinsJobServer is an implementation of IJobServer.
type JenkinsJobServer struct {
	jenkins *gojenkins.Jenkins
	jobs    *JobList
}

// Connect establishes connection to the JenkinsJobServer.
func (js *JenkinsJobServer) Connect(url, user, password string) error {
	js.jenkins = gojenkins.CreateJenkins(nil, url, user, password)
	_, err := js.jenkins.Init()
	if err != nil {
		return err
	}
	js.jobs = &JobList{}
	js.jobs.Clear()
	err = js.Load()
	if err != nil {
		return err
	}
	return nil
}

// Load queries the job server for all the data.
func (js *JenkinsJobServer) Load() error {
	jobs, err := js.jenkins.GetAllJobs()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		js.jobs.AddJob(&Job{
			jenkinsJob: job,
			client:     js.jenkins,
		})
	}
	return nil
}

// GetJobs returns an IJobList implementation.
func (js *JenkinsJobServer) GetJobs() IJobList {
	return js.jobs
}

// GetJob returns a specific IJob by name.
func (js *JenkinsJobServer) GetJob(jobName string) IJob {
	return js.jobs.GetJob(jobName)
}
