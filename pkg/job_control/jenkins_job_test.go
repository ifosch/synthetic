package jobcontrol

import (
	"testing"

	"github.com/bndr/gojenkins"
)

func TestJob(t *testing.T) {
	j := &Job{
		jenkinsJob: &gojenkins.Job{
			Raw: &gojenkins.JobResponse{
				Name:        "myjob",
				Description: "myjob does something",
			},
		},
	}

	name := j.Name()

	description := j.Description()

	describe := j.Describe()

	if name != "myjob" {
		t.Logf("Wrong job name '%v' should be 'myjob'", name)
		t.Fail()
	}
	if description != "myjob does something" {
		t.Logf("Wrong job description '%v' should be 'myjob does something'", description)
		t.Fail()
	}
	if describe != "`myjob`: myjob does something\nParameters:\n" {
		t.Logf("Wrong job describe '%v' should be '`myjob`: myjob does something\nParameters\n'", describe)
		t.Fail()
	}
}
