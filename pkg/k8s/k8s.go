package k8s

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/ifosch/synthetic/pkg/slack"
	"github.com/ifosch/synthetic/pkg/synthetic"
)

// GetClusters loads default kubeconfig and gets the list of cluster
// defined.
func GetClusters() ([]string, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

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

func listClusters(msg synthetic.Message) {
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

// Register registers all the kubernetes operations as bot commands.
func Register(client *slack.Chat) {
	client.RegisterMessageProcessor(slack.NewMessageProcessor("github.com/ifosch/synthetic/pkg/k8s.listClusters", slack.Exactly(slack.Mentioned(listClusters), "list clusters")))
}
