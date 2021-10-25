package k8s

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

func getConfig(context string) clientcmd.ClientConfig {
	configOverrides := &clientcmd.ConfigOverrides{}

	if context != "" {
		configOverrides.CurrentContext = context
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		configOverrides,
	)
}

var getClient = func(cluster string) (kubernetes.Interface, error) {
	kubeConfig := getConfig(cluster)

	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetPods returns a list of pods for the specified `cluster`. If
// `cluster` is an empty string, then it lists all pods in all
// clusters known by the bot.
func GetPods(cluster, namespace string) ([]v1.Pod, error) {
	client, err := getClient(cluster)
	if err != nil {
		return nil, err
	}

	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

// ListPods returns a list of pods
func ListPods(msg synthetic.Message) {
	command := strings.Split(msg.ClearMention(), " ")
	cluster := ""
	namespace := ""
	if len(command) >= 3 {
		cluster = command[2]
		if len(command) == 4 {
			namespace = command[3]
		}
	}

	pods, err := GetPods(cluster, namespace)
	if err != nil {
		msg.Reply(err.Error(), msg.Thread())
		return
	}

	response := ""
	for _, pod := range pods {
		response = fmt.Sprintf("%s- %s\n", response, pod.Name)
	}
	msg.Reply(response, msg.Thread())
}

// GetClusters loads default kubeconfig and gets the list of cluster
// defined.
func GetClusters() ([]string, error) {
	kubeConfig := getConfig("")

	clusters := []string{}
	cfg, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, err
	}
	for cluster := range cfg.Clusters {
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// ListClusters returns a list of clusters available in the supplied
// kubeconfig
func ListClusters(msg synthetic.Message) {
	clusters, err := GetClusters()
	if err != nil {
		msg.Reply(err.Error(), msg.Thread())
		return
	}
	if len(clusters) < 1 {
		msg.Reply("I know of no kubernetes clusters. Checkout my kubeconfig.", msg.Thread())
		return
	}
	response := "I know of the following clusters:\n"
	for _, cluster := range clusters {
		response = fmt.Sprintf("%s- %s\n", response, cluster)
	}
	msg.Reply(response, msg.Thread())
}
