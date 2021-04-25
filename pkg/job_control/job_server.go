package jobcontrol

// IJobServer is an interface to a job server.
type IJobServer interface {
	Connect(string, string, string) error
	Load() error
	GetJobs() IJobList
	GetJob(string) IJob
}
