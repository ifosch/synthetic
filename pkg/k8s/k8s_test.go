package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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

func TestGetPods(t *testing.T) {
	tcs := []struct {
		namespace string
		podList   []*v1.Pod
	}{
		{
			namespace: "",
			podList: []*v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod1",
						Namespace: "default",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod2",
						Namespace: "kube-system",
					},
				},
			},
		},
		{
			namespace: "kube-system",
			podList: []*v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod2",
						Namespace: "kube-system",
					},
				},
			},
		},
	}

	for _, test := range tcs {
		clientSet := fake.NewSimpleClientset()
		getClient = func(cluster string) (kubernetes.Interface, error) {
			for _, pod := range test.podList {
				clientSet.CoreV1().Pods(pod.Namespace).Create(context.Background(), pod, metav1.CreateOptions{})
			}
			return clientSet, nil
		}

		pods, err := GetPods("", test.namespace)
		if err != nil {
			panic(err)
		}

		if len(pods) != len(test.podList) {
			t.Errorf("Wrong number of pods %d, expected %d", len(pods), len(test.podList))
		}
		for i, pod := range pods {
			if pod.Name != test.podList[i].Name {
				t.Errorf("Wrong pod name %s, expected %s", pods[0].Name, test.podList[i].Name)
			}
		}
	}
}

func TestGetConfig(t *testing.T) {
	tcs := []struct {
		kubeconfig    string
		chosenContext string
		expectedHost  string
	}{
		{
			kubeconfig:    getKubeCfgFixture("kubeconfig.yaml"),
			chosenContext: "",
			expectedHost:  "https://cluster1.example.com",
		},
		{
			kubeconfig:    getKubeCfgFixture("kubeconfig.yaml"),
			chosenContext: "cluster3.example.com",
			expectedHost:  "https://cluster3.example.com",
		},
	}

	for _, test := range tcs {
		os.Setenv("KUBECONFIG", test.kubeconfig)

		kubeConfig := getConfig(test.chosenContext)

		clientConfig, err := kubeConfig.ClientConfig()
		if err != nil {
			panic(err)
		}

		if clientConfig.Host != test.expectedHost {
			t.Fatalf("Incorrect host: got %s, expected %s", clientConfig.Host, test.expectedHost)
		}
	}
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
