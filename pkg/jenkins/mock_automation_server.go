package jenkins

import (
	"sort"
	"text/template"

	jk "github.com/ifosch/jk/pkg/jenkins"
)

// MockAutomationServer is a mocking AutomationServer for testing.
type MockAutomationServer struct {
	jobs     map[string]string
	jobNames []string
}

// NewMockAutomationServer returns a mocking AutomationServer.
func NewMockAutomationServer(jobs map[string]string) *MockAutomationServer {
	jobNames := []string{}
	for jobName := range jobs {
		jobNames = append(jobNames, jobName)
	}
	sort.Strings(jobNames)
	return &MockAutomationServer{
		jobs:     jobs,
		jobNames: jobNames,
	}
}

// List is a goroutine that puts a Message for each job defined in the
// mocking AutomationServer.
func (mas *MockAutomationServer) List(out chan jk.Message) {
	for _, j := range mas.jobNames {
		msg := jk.Message{
			Message: j,
			Error:   false,
			Done:    false,
		}
		out <- msg
	}
	out <- jk.Message{
		Message: "",
		Error:   false,
		Done:    true,
	}
}

// Describe is a goroutine that puts a Message with the specified job
// description.
func (mas *MockAutomationServer) Describe(name string, t *template.Template, out chan jk.Message) {
	if _, ok := mas.jobs[name]; !ok {
		out <- jk.Message{
			Message: "Non-existing job name",
			Error:   true,
			Done:    true,
		}
		return
	}
	out <- jk.Message{
		Message: mas.jobs[name],
		Error:   false,
		Done:    true,
	}
}

// Build is a gorouting that puts a Message with the updates of a
// running job.
func (mas *MockAutomationServer) Build(name string, args map[string]string, out chan jk.Message) {
	out <- jk.Message{
		Message: "Build queued",
		Error:   false,
		Done:    false,
	}
	out <- jk.Message{
		Message: "Build started",
		Error:   false,
		Done:    false,
	}
	out <- jk.Message{
		Message: "Build finished",
		Error:   false,
		Done:    true,
	}
}
