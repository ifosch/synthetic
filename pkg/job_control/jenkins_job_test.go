package jobcontrol

import (
	"testing"

	"github.com/bndr/gojenkins"
)

type jenkinsJobTC struct {
	job              *gojenkins.Job
	expectedDescribe string
}

func TestJob(t *testing.T) {
	tc := jenkinsJobTC{
		job: &gojenkins.Job{
			Raw: &gojenkins.JobResponse{
				Name:        "myjob",
				Description: "myjob does something",
			},
		},
		expectedDescribe: "`myjob`: myjob does something\nParameters:\n",
	}
	j := &Job{
		jenkinsJob: tc.job,
	}

	name := j.Name()

	description := j.Description()

	describe := j.Describe()

	if name != tc.job.Raw.Name {
		t.Errorf("Wrong job name '%v' should be '%v'", name, tc.job.Raw.Name)
	}
	if description != tc.job.Raw.Description {
		t.Errorf("Wrong job description '%v' should be '%v'", description, tc.job.Raw.Description)
	}
	if describe != tc.expectedDescribe {
		t.Errorf("Wrong job describe '%v' should be '%v'", describe, tc.expectedDescribe)
	}
}
