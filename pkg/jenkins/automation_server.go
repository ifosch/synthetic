package jenkins

import (
	jk "github.com/ifosch/jk/pkg/jenkins"
	"text/template"
)

// AutomationServer is an interface to represent servers to automate
// tasks, like Jenkins.
type AutomationServer interface {
	List(chan jk.Message)
	Describe(string, *template.Template, chan jk.Message)
	Build(string, map[string]string, chan jk.Message)
}
