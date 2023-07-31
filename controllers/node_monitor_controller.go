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
	monitoring.ControllerStatus("NodeMonitorController", true)
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

func (n *NodeMonitorController) logNodeDetails(event string, node *v1.Node, status v1.ConditionStatus) {
	n.Logger.Infof("Node %s: %s/%s", event, node.Name, status)
	n.Logger.Debug("Node Details: ", node.Status, node.Status.Conditions, node.Annotations, node.Labels, node.Finalizers, node.Spec.Taints, node.Status.Capacity, node.Status.Allocatable, node.Status.Addresses)
}

func (n *NodeMonitorController) onAdd(obj interface{}) {
	node := obj.(*v1.Node)
	status, _ := getNodeReadyStatus(node)
	n.logNodeDetails("added", node, status)
	summary := NewNodeSummary(node)
	monitoring.UpdateNodeMetrics(nil, &summary)
}

func (n *NodeMonitorController) onUpdate(oldObj, newObj interface{}) {
	oldNode, ok1 := oldObj.(*v1.Node)
	newNode, ok2 := newObj.(*v1.Node)
	if !ok1 || !ok2 {
		n.Logger.Error("Failed to convert obj to Node")
		return
	}
	status, _ := getNodeReadyStatus(newNode)
	n.logNodeDetails("updated", newNode, status)
	oldSummary := NewNodeSummary(oldNode)
	newSummary := NewNodeSummary(newNode)
	monitoring.UpdateNodeMetrics(&oldSummary, &newSummary)
}

func (n *NodeMonitorController) onDelete(obj interface{}) {
	node, ok := obj.(*v1.Node)
	if !ok {
		fmt.Println("Delete event with incorrect type:", obj)
		return
	}
	status, _ := getNodeReadyStatus(node)
	n.logNodeDetails("deleted", node, status)
	summary := NewNodeSummary(node)
	monitoring.UpdateNodeMetrics(&summary, nil)
}

// NewNodeSummary creates a NodeSummary from a given node
func NewNodeSummary(node *v1.Node) monitoring.NodeSummary {
	readyStatus, _ := getNodeReadyStatus(node)
	return monitoring.NodeSummary{
		Name:        node.Name,
		ReadyStatus: readyStatus,
		Conditions:  node.Status.Conditions,
		Annotations: node.Annotations,
		Labels:      node.Labels,
		Taints:      node.Spec.Taints,
		Capacity:    node.Status.Capacity,
		Allocatable: node.Status.Allocatable,
	}
}

// getNodeReadyStatus is a helper function to extract the ready status from a node
func getNodeReadyStatus(node *v1.Node) (v1.ConditionStatus, string) {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			if condition.Status == v1.ConditionTrue {
				return condition.Status, "Ready"
			} else {
				return condition.Status, "NotReady"
			}
		}
	}
	return v1.ConditionUnknown, "Unknown"
}
