package jobcontrol

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/bndr/gojenkins"
)

// Job is an implementation of an IJob.
type Job struct {
	client           *gojenkins.Jenkins
	jenkinsJob       *gojenkins.Job
	describeTemplate *template.Template
}

// Name returns the job name.
func (j *Job) Name() string {
	return j.jenkinsJob.GetName()
}

// Description return the job's description.
func (j *Job) Description() string {
	return j.jenkinsJob.GetDescription()
}

// Run runs the Job.
func (j *Job) Run(args map[string]string, out chan Update) {
	number, err := j.client.BuildJob(j.Name(), args)
	if err != nil {
		update(out, fmt.Sprintf("Job Invoke error %v", err), "boom", true)
		return
	}
	task, err := j.client.GetQueueItem(number)
	if err != nil {
		update(out, fmt.Sprintf("Task get error %v", err), "boom", true)
		return
	}
	update(out, fmt.Sprintf("Execution for job `%v` was queued", j.Name()), "stopwatch", false)
	buildID := task.Raw.Executable.Number
	for {
		if buildID != 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
		task.Poll()
		buildID = task.Raw.Executable.Number
	}
	build, err := j.client.GetBuild(j.Name(), buildID)
	if err != nil {
		update(out, fmt.Sprintf("Queue item get error %v", err), "boom", true)
		return
	}
	update(out, fmt.Sprintf("Building `%v` with parameters `%v` (%v)", j.Name(), args, build.GetUrl()), "gear", false)
	for {
		if !build.Raw.Building {
			break
		}
		time.Sleep(100 * time.Millisecond)
		_, err = build.Poll()
		if err != nil {
			update(out, fmt.Sprintf("Error polling build %v", err), "boom", true)
			return
		}
	}
	update(out, fmt.Sprintf("Job `%v` completed", j.Name()), "heavy_check_mark", true)
}

// Describe describes the Job.
func (j *Job) Describe() string {
	var err error
	if j.describeTemplate == nil {
		j.describeTemplate, err = DescribeTemplate()
		if err != nil {
			return fmt.Sprintf("Template parsing error: %s", err)
		}
	}
	msg := &bytes.Buffer{}
	err = j.describeTemplate.Execute(msg, j.jenkinsJob)
	if err != nil {
		return fmt.Sprintf("Template execution error: %s", err)
	}
	return msg.String()
}

// DescribeTemplate returns the default template for Describe.
func DescribeTemplate() (t *template.Template, err error) {
	tmpl := "`{{ .Raw.Name }}`: "
	tmpl += "{{ .Raw.Description }}\n"
	tmpl += "Parameters:\n"
	tmpl += "{{ with .Raw.Property }}"
	tmpl += "{{ range . }}"
	tmpl += "{{ with .ParameterDefinitions }}"
	tmpl += "{{ range . }}"
	tmpl += "- *{{ .Name }}* ({{ .Type }}): "
	tmpl += "{{ trim .Description }} "
	tmpl += "(Default: `{{ .DefaultParameterValue.Value }}`)\n"
	tmpl += "{{ end }}"
	tmpl += "{{ end }}"
	tmpl += "{{ end }}"
	tmpl += "{{ end }}"
	t, err = template.New("Job").Funcs(template.FuncMap{
		"trim": trim,
	}).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("Template Creation error %v", err)
	}
	return t, nil
}

// trim is a helper function for templates to remove surrounding
// newlines.
func trim(in string) string {
	return strings.TrimPrefix(
		strings.TrimSuffix(in, "\n"),
		"\n",
	)
}

func update(out chan Update, reply, reaction string, done bool) {
	out <- Update{
		Msg:      reply,
		Reaction: reaction,
		Done:     done,
	}
}
