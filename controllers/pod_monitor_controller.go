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

func NewPodMonitorController(clientset KubernetesClientset, logger *logrus.Logger) *PodMonitorController {
	return &PodMonitorController{
		Clientset:   clientset,
		Logger:      logger.WithField("controller", "PodMonitorController"),
		podStatuses: make(map[string]int),
	}
}

func (n *PodMonitorController) Start() error {
	monitoring.ControllerStatus("PodMonitorController", true)
	factory := informers.NewSharedInformerFactory(n.Clientset, time.Minute)
	PodInformer := factory.Core().V1().Pods().Informer()

	PodInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    n.onAdd,
		UpdateFunc: n.onUpdate,
		DeleteFunc: n.onDelete,
	})

	n.stopCh = make(chan struct{})
	go PodInformer.Run(n.stopCh)

	if !cache.WaitForCacheSync(n.stopCh, PodInformer.HasSynced) {
		logrus.Error("Timed out waiting for caches to sync")
		return fmt.Errorf("timed out waiting for caches to sync")
	}

	logrus.Debug("Pod informer synced successfully")
	return nil
}

func (n *PodMonitorController) Stop() error {
	monitoring.ControllerStatus("PodMonitorController", false)
	close(n.stopCh)
	return nil
}

func (n *PodMonitorController) logPodDetails(Pod *v1.Pod, event string) {
	n.Logger.Infof("Pod %s: %s (status: %s)", event, Pod.Name, Pod.Status.Phase)
	n.Logger.Debug("Pod Details: ", Pod.Status, Pod.Status.Conditions, Pod.Annotations, Pod.Labels, Pod.Finalizers, Pod.Status.PodIP)
}

func (n *PodMonitorController) onAdd(obj interface{}) {
	Pod := obj.(*v1.Pod)
	n.logPodDetails(Pod, "added")
	monitoring.PodCountChange("inc")
	status := string(Pod.Status.Phase)
	n.podStatuses[status]++
	monitoring.UpdatePodStatus(status, float64(n.podStatuses[status]))
}

func (n *PodMonitorController) onUpdate(oldObj, newObj interface{}) {
	oldPod := oldObj.(*v1.Pod)
	newPod := newObj.(*v1.Pod)
	n.logPodDetails(newPod, "updated")

	oldStatus := string(oldPod.Status.Phase)
	newStatus := string(newPod.Status.Phase)
	n.podStatuses[oldStatus]--
	n.podStatuses[newStatus]++
	monitoring.UpdatePodStatus(oldStatus, float64(n.podStatuses[oldStatus]))
	monitoring.UpdatePodStatus(newStatus, float64(n.podStatuses[newStatus]))

	oldSummary := monitoring.CreatePodSummary(oldPod)
	newSummary := monitoring.CreatePodSummary(newPod)

	for _, container := range oldSummary.Containers {
		if container.IsCrashLooping {
			n.CrashLoopingCount--
		}
	}

	for _, container := range newSummary.Containers {
		if container.IsCrashLooping {
			n.CrashLoopingCount++
		}
	}

	monitoring.UpdatePodStatus("CrashLoopBackOff", float64(n.CrashLoopingCount))
}

func (n *PodMonitorController) onDelete(obj interface{}) {
	Pod, ok := obj.(*v1.Pod)
	if !ok {
		fmt.Println("Delete event with incorrect type:", obj)
		return
	}
	n.logPodDetails(Pod, "deleted")
	monitoring.PodCountChange("dec")
	status := string(Pod.Status.Phase)
	n.podStatuses[status]--
	monitoring.UpdatePodStatus(status, float64(n.podStatuses[status]))
}
