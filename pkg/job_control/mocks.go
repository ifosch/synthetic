package jobcontrol

import (
	"fmt"
	"os"
)

// MockJobServer is a mocking JobServer for testing.
type MockJobServer struct {
	jobs         IJobList
	originalJobs map[string]string
}

// NewMockJobServer returns a mocking JobServer.
func NewMockJobServer(jobs map[string]string) *MockJobServer {
	jobServer := &MockJobServer{
		jobs:         &JobList{},
		originalJobs: jobs,
	}
	jobServer.Load()
	return jobServer
}

// Connect mocks JobServer.Connect method.
func (mjs *MockJobServer) Connect(string, string, string) error {
	return nil
}

// Load mocks JobServer.Load method.
func (mjs *MockJobServer) Load() error {
	for jobName := range mjs.originalJobs {
		mjs.jobs.AddJob(&MockJob{
			name:        jobName,
			description: mjs.originalJobs[jobName],
		})

	}
	return nil
}

// GetJobs mocks JobServer.GetJobs method.
func (mjs *MockJobServer) GetJobs() IJobList {
	return mjs.jobs
}

// GetJob mocks JobServer.GetJob method.
func (mjs *MockJobServer) GetJob(jobName string) IJob {
	return mjs.GetJobs().GetJob(jobName)
}

// MockJob mocks a Job.
type MockJob struct {
	name        string
	description string
}

// Name mocks Job.Name method.
func (j *MockJob) Name() string {
	return j.name
}

// Description mocks Job.Description method.
func (j *MockJob) Description() string {
	return j.description
}

// Run mocks Job.Run method.
func (j *MockJob) Run(args map[string]string, out chan Update) {
	out <- Update{
		Msg: fmt.Sprintf(
			"Execution for job `%s` was queued",
			j.name,
		),
		Reaction: "stopwatch",
		Done:     false,
	}
	out <- Update{
		Msg: fmt.Sprintf(
			"Building `%s` with parameters `%s` (%s/job/%s)",
			j.name,
			args,
			os.Getenv("JENKINS_URL"),
			j.name,
		),
		Reaction: "gear",
		Done:     false,
	}
	out <- Update{
		Msg: fmt.Sprintf(
			"Job %s completed",
			j.name,
		),
		Reaction: "heavy_check_mark",
		Done:     true,
	}
}

// Describe mocks Job.Describe method.
func (j *MockJob) Describe() string {
	return j.Description()
}
