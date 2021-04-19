package jobcontrol

// IJob is an interface to a job.
type IJob interface {
	Name() string
	Description() string
	Run(map[string]string, chan Update)
	Describe() string
}
