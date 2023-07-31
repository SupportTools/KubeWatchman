package k8s

import (
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type ClusterConnection interface {
	RESTClient() rest.Interface
	CoreV1() v1.CoreV1Interface
}

type ClusterConnectionFactory interface {
	InClusterConfig() (*rest.Config, error)
	NewForConfig(*rest.Config) (*kubernetes.Clientset, error)
}
