package jobcontrol

import (
	"testing"

	"github.com/bndr/gojenkins"
)

type jenkinsJobTC struct {
	input            string
	job              *gojenkins.Job
	expectedDescribe string
}

func TestJob(t *testing.T) {
	tcs := []jenkinsJobTC{
		{
			input: "Simple job without params",
			job: &gojenkins.Job{
				Raw: &gojenkins.JobResponse{
					Name:        "myjob",
					Description: "myjob does something",
				},
			},
			expectedDescribe: "`myjob`" + `: myjob does something
Parameters:
`,
		},
		{
			input: "Simple job with one param",
			job: &gojenkins.Job{
				Raw: &gojenkins.JobResponse{
					Name:        "myjob",
					Description: "myjob does something",
					Property: []struct {
						ParameterDefinitions []gojenkins.ParameterDefinition `json:"parameterDefinitions"`
					}{
						{
							ParameterDefinitions: []gojenkins.ParameterDefinition{
								{
									Description: "myParam helps parametrize myjob",
									Name:        "myParam",
									Type:        "StringParameterDefinition",
									DefaultParameterValue: struct {
										Name  string      `json:"name"`
										Value interface{} `json:"value"`
									}{
										"myParam",
										"all",
									},
								},
							},
						},
					},
				},
			},
			expectedDescribe: "`myjob`" + `: myjob does something
Parameters:
- *myParam* (StringParameterDefinition): myParam helps parametrize myjob (Default: ` + "`all`)\n",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			j := &Job{
				jenkinsJob: tc.job,
			}

			name := j.Name()
			if name != tc.job.Raw.Name {
				t.Errorf("Wrong job name '%v' should be '%v'", name, tc.job.Raw.Name)
			}

			description := j.Description()
			if description != tc.job.Raw.Description {
				t.Errorf("Wrong job description '%v' should be '%v'", description, tc.job.Raw.Description)
			}

			describe := j.Describe()
			if describe != tc.expectedDescribe {
				t.Errorf("Wrong job describe '%v' should be '%v'", describe, tc.expectedDescribe)
			}
		})
	}
}
