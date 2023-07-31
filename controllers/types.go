package controllers

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type KubernetesClientset interface {
	kubernetes.Interface
}

type NodeMonitorController struct {
	Clientset KubernetesClientset
	Logger    *logrus.Entry
	stopCh    chan struct{}
}
