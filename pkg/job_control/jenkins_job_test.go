package jobcontrol

import (
	"fmt"
	"testing"

	"github.com/bndr/gojenkins"
)

type jenkinsJobTC struct {
	name        string
	description string
	params      []struct {
		Name         string
		Type         string
		Description  string
		DefaultValue string
	}
}

func (tc jenkinsJobTC) Describe() string {
	expectedDescription := fmt.Sprintf(
		"`%s`: %s\nParameters:\n",
		tc.name,
		tc.description,
	)

	for _, param := range tc.params {
		expectedDescription += fmt.Sprintf(
			"- *%s* (%s): %s (Default: `%s`)\n",
			param.Name,
			param.Type,
			param.Description,
			param.DefaultValue,
		)
	}

	return expectedDescription
}

func (tc jenkinsJobTC) getJenkinsJob() *gojenkins.Job {
	parameterDefinitions := []gojenkins.ParameterDefinition{}
	for _, param := range tc.params {
		parameterDefinition := gojenkins.ParameterDefinition{
			Description: param.Description,
			Name:        param.Name,
			Type:        param.Type,
			DefaultParameterValue: struct {
				Name  string      `json:"name"`
				Value interface{} `json:"value"`
			}{
				Name:  "",
				Value: param.DefaultValue,
			},
		}
		parameterDefinitions = append(parameterDefinitions, parameterDefinition)
	}
	return &gojenkins.Job{
		Raw: &gojenkins.JobResponse{
			Name:        tc.name,
			Description: tc.description,
			Property: []struct {
				ParameterDefinitions []gojenkins.ParameterDefinition `json:"parameterDefinitions"`
			}{
				{
					ParameterDefinitions: parameterDefinitions,
				},
			},
		},
	}
}

func TestJob(t *testing.T) {
	tcs := map[string]jenkinsJobTC{
		"Simple job without params": {
			name:        "myjob",
			description: "myjob does something",
			params: []struct {
				Name         string
				Type         string
				Description  string
				DefaultValue string
			}{},
		},
		"Simple job with one param": {
			name:        "myjob",
			description: "myjob does something",
			params: []struct {
				Name         string
				Type         string
				Description  string
				DefaultValue string
			}{
				{
					Name:         "myParam",
					Type:         "StringParameterDefinition",
					Description:  "myParam helps parametrize myjob",
					DefaultValue: "all",
				},
			},
		},
	}

	for testID, tc := range tcs {
		t.Run(testID, func(t *testing.T) {
			j := &Job{
				jenkinsJob: tc.getJenkinsJob(),
			}

			name := j.Name()
			if name != tc.name {
				t.Errorf("Wrong job name '%v' should be '%v'", name, tc.name)
			}

			description := j.Description()
			if description != tc.description {
				t.Errorf("Wrong job description '%v' should be '%v'", description, tc.description)
			}

			describe := j.Describe()
			if describe != tc.Describe() {
				t.Errorf("Wrong job describe '%v' should be '%v'", describe, tc.Describe())
			}
		})
	}
}
