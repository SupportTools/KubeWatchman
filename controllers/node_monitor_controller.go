package controllers

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/supporttools/KubeWatchman/monitoring"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func NewNodeMonitorController(clientset KubernetesClientset, logger *logrus.Logger) *NodeMonitorController {
	return &NodeMonitorController{
		Clientset: clientset,
		Logger:    logger.WithField("controller", "NodeMonitorController"),
	}
}

func (n *NodeMonitorController) Start() error {
	monitoring.ControllerStatus("nodeMonitorController", true)
	// Create a shared informer factory
	factory := informers.NewSharedInformerFactory(n.Clientset, time.Minute)
	logrus.Debug("Created shared informer factory")

	// Create a Node informer
	nodeInformer := factory.Core().V1().Nodes().Informer()
	logrus.Debug("Created Node informer")

	// Set up the event handlers for changes
	//nolint:errcheck
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    n.onAdd,
		UpdateFunc: n.onUpdate,
		DeleteFunc: n.onDelete,
	})

	// Start the informer
	n.stopCh = make(chan struct{})
	go nodeInformer.Run(n.stopCh)

	// Wait for the informer's cache to sync
	if !cache.WaitForCacheSync(n.stopCh, nodeInformer.HasSynced) {
		logrus.Error("Timed out waiting for caches to sync")
		return fmt.Errorf("timed out waiting for caches to sync")
	}

	logrus.Debug("Node informer synced successfully")
	return nil
}

func (n *NodeMonitorController) Stop() error {
	monitoring.ControllerStatus("nodeMonitorController", false)
	close(n.stopCh)
	return nil
}

func (n *NodeMonitorController) onAdd(obj interface{}) {
	node := obj.(*v1.Node)
	_, msg := getNodeReadyStatus(node)
	n.Logger.Info("Node added:", node.Name, msg)
	n.Logger.Debug("Node Status: ", node.Status)
	n.Logger.Debug("Node Conditions: ", node.Status.Conditions)
	n.Logger.Debug("Node Annotations: ", node.Annotations)
	n.Logger.Debug("Node Labels: ", node.Labels)
	n.Logger.Debug("Node Finalizers: ", node.Finalizers)
	n.Logger.Debug("Node Taints: ", node.Spec.Taints)
	n.Logger.Debug("Node Capacity: ", node.Status.Capacity)
	n.Logger.Debug("Node Allocatable: ", node.Status.Allocatable)
	n.Logger.Debug("Node Addresses: ", node.Status.Addresses)
}

func (n *NodeMonitorController) onUpdate(oldObj, newObj interface{}) {
	node := newObj.(*v1.Node)
	_, msg := getNodeReadyStatus(node)
	n.Logger.Info("Node updated:", node.Name, msg)
	n.Logger.Debug("Node Status: ", node.Status)
	n.Logger.Debug("Node Conditions: ", node.Status.Conditions)
	n.Logger.Debug("Node Annotations: ", node.Annotations)
	n.Logger.Debug("Node Labels: ", node.Labels)
	n.Logger.Debug("Node Finalizers: ", node.Finalizers)
	n.Logger.Debug("Node Taints: ", node.Spec.Taints)
	n.Logger.Debug("Node Capacity: ", node.Status.Capacity)
	n.Logger.Debug("Node Allocatable: ", node.Status.Allocatable)
	n.Logger.Debug("Node Addresses: ", node.Status.Addresses)
}

func (n *NodeMonitorController) onDelete(obj interface{}) {
	node, ok := obj.(*v1.Node)
	if !ok {
		fmt.Println("Delete event with incorrect type:", obj)
		return
	}
	status, msg := getNodeReadyStatus(node)
	n.Logger.Info("Node deleted:", node.Name, msg)
	n.Logger.Debug("Node Status: ", status)
	n.Logger.Debug("Node Conditions: ", node.Status.Conditions)
	n.Logger.Debug("Node Annotations: ", node.Annotations)
	n.Logger.Debug("Node Labels: ", node.Labels)
	n.Logger.Debug("Node Finalizers: ", node.Finalizers)
	n.Logger.Debug("Node Taints: ", node.Spec.Taints)
	n.Logger.Debug("Node Addresses: ", node.Status.Addresses)
	n.Logger.Debug("Node Capacity: ", node.Status.Capacity)
	n.Logger.Debug("Node Allocatable: ", node.Status.Allocatable)
	n.Logger.Debug("Node Conditions: ", node.Status.Conditions)
}

func getNodeReadyStatus(node *v1.Node) (v1.ConditionStatus, string) {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return condition.Status, "Ready"
		}
	}
	return v1.ConditionUnknown, "Unknown"
}
