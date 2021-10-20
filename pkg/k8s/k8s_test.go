package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func stringIn(item string, items []string) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}
	return false
}

func compareStringLists(a, b []string) error {
	for _, item := range a {
		if !stringIn(item, b) {
			return fmt.Errorf("item %s was not expected", item)
		}
	}
	for _, item := range b {
		if !stringIn(item, a) {
			return fmt.Errorf("item %s is missing", item)
		}
	}

	return nil
}

func getKubeCfgFixture(filename string) string {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(curDir, fmt.Sprintf("../../tests/fixtures/%s", filename))
}

func TestListClusters(t *testing.T) {
	tcs := []struct {
		kubeconfig     string
		expectedErrMsg string
	}{
		{
			kubeconfig:     getKubeCfgFixture("kubeconfig.yaml"),
			expectedErrMsg: "",
		},
		{
			kubeconfig:     getKubeCfgFixture("kubeconfig_one_missing.yaml"),
			expectedErrMsg: "item cluster3.example.com is missing",
		},
		{
			kubeconfig:     getKubeCfgFixture("kubeconfig_one_unexpected.yaml"),
			expectedErrMsg: "item unexpected.example.com was not expected",
		},
	}

	for _, test := range tcs {
		os.Setenv("KUBECONFIG", test.kubeconfig)
		expectedClusters := []string{
			"cluster1.example.com",
			"cluster2.example.com",
			"cluster3.example.com",
		}

		clusters, err := GetClusters()
		if err != nil {
			if err.Error() != test.expectedErrMsg {
				t.Fatalf("%s", err)
			}
		}

		err = compareStringLists(clusters, expectedClusters)
		if err != nil {
			if err.Error() != test.expectedErrMsg {
				t.Fatalf("%s", err)
			}
		}
	}
}
