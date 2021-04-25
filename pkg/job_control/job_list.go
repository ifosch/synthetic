package jobcontrol

import (
	"fmt"
)

// IJobList is an interface to a collection of jobs.
type IJobList interface {
	AddJob(IJob)
	Len() int
	Clear()
	GetJob(string) IJob
	List() string
}

// JobList is an implementation of an IJobList.
type JobList struct {
	jobs []IJob
}

// AddJob adds a new IJob to JobList.
func (jl *JobList) AddJob(job IJob) {
	jl.jobs = append(jl.jobs, job)
}

// Len returns the length of the job list.
func (jl *JobList) Len() int {
	return len(jl.jobs)
}

// Clear resets the object.
func (jl *JobList) Clear() {
	jl.jobs = []IJob{}
}

// GetJob retrieves a job identified by name.
func (jl *JobList) GetJob(name string) IJob {
	for _, job := range jl.jobs {
		if job.Name() == name {
			return job
		}
	}
	return nil
}

// List returns a string with a list of job names.
func (jl *JobList) List() string {
	result := ""
	for _, job := range jl.jobs {
		result = fmt.Sprintf("%s%s\n", result, job.Name())
	}
	return result
}
