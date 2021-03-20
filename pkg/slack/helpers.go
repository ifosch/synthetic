package slack

import (
	"reflect"
	"runtime"
)

// getProcessorName ...
func getProcessorName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
