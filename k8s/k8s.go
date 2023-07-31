package k8s

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func CreateClusterConnection(logger *logrus.Logger, factory ClusterConnectionFactory) (*kubernetes.Clientset, error) {
	// Use the factory to create the in-cluster config and clientset
	config, err := factory.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := factory.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func CheckClusterConnection(clientset ClusterConnection, logger *logrus.Logger) error {
	if clientset == nil {
		return fmt.Errorf("failed to create Kubernetes clientset")
	}

	ctx := context.TODO()
	if err := clientset.CoreV1().RESTClient().Get().Do(ctx).Error(); err != nil {
		return fmt.Errorf("failed to connect to the cluster: %v", err)
	}

	logger.Debug("Successfully connected to the cluster")
	return nil
}
