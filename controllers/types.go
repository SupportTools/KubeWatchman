package controllers

import (
	"github.com/sirupsen/logrus"
	"github.com/supporttools/KubeWatchman/monitoring"
	"k8s.io/client-go/kubernetes"
)

type KubernetesClientset interface {
	kubernetes.Interface
}

type NodeMonitorController struct {
	Clientset   KubernetesClientset
	Logger      *logrus.Entry
	stopCh      chan struct{}
	NodeSummary map[string]*monitoring.NodeSummary
}

type PodMonitorController struct {
	Clientset         KubernetesClientset
	Logger            *logrus.Entry
	stopCh            chan struct{}
	PodSummary        map[string]*monitoring.PodSummary
	podStatuses       map[string]int
	CrashLoopingCount int
}
