package monitoring

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/api/core/v1"
)

var (
	controllersUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubewatchman_controllers_up",
			Help: "Status of the controllers",
		},
		[]string{"controller"},
	)
	nodeCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kubewatchman_node_count",
			Help: "Number of nodes in the cluster",
		},
	)
	nodeStatusReady = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kubewatchman_node_status_ready",
			Help: "Count of the nodes in ready state",
		},
	)
	nodeStatusNotready = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kubewatchman_node_status_notready",
			Help: "Count of the nodes in not ready state",
		},
	)
	podCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kubewatchman_pod_count",
			Help: "Number of pods in the cluster",
		},
	)
	podStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubewatchman_pod_status",
			Help: "Status of the pods",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(controllersUp)
	prometheus.MustRegister(nodeCount)
	prometheus.MustRegister(nodeStatusReady)
	prometheus.MustRegister(nodeStatusNotready)
	prometheus.MustRegister(podCount)
	prometheus.MustRegister(podStatus)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func ControllerStatus(controller string, status bool) {
	if status {
		controllersUp.WithLabelValues(controller).Set(1)
	} else {
		controllersUp.WithLabelValues(controller).Set(0)
	}
}

// Node Section
func NodeStatus(input string) {
	switch input {
	case "inc":
		nodeCount.Inc()
	case "dec":
		nodeCount.Dec()
	}
}

func NodeCountChange(input string) {
	switch input {
	case "inc":
		nodeCount.Inc()
	case "dec":
		nodeCount.Dec()
	}
}

func UpdateNodeMetrics(oldNodeSummary *NodeSummary, newNodeSummary *NodeSummary) {
	// If oldNodeSummary is not nil, it means we have to decrement the old status
	if oldNodeSummary != nil {
		if oldNodeSummary.ReadyStatus == v1.ConditionTrue {
			nodeStatusReady.Dec()
		} else {
			nodeStatusNotready.Dec()
		}
	}

	// Increment the new status
	if newNodeSummary.ReadyStatus == v1.ConditionTrue {
		nodeStatusReady.Inc()
	} else {
		nodeStatusNotready.Inc()
	}
}

// Pod Section
func PodCountChange(input string) {
	switch input {
	case "inc":
		podCount.Inc()
	case "dec":
		podCount.Dec()
	}
}

func UpdatePodStatus(status string, value float64) {
	podStatus.WithLabelValues(status).Set(value)
}

func CreatePodSummary(pod *v1.Pod) *PodSummary {
	containers := make([]ContainerSummary, len(pod.Status.ContainerStatuses))
	for i, containerStatus := range pod.Status.ContainerStatuses {
		containers[i] = ContainerSummary{
			Name:           containerStatus.Name,
			IsCrashLooping: containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Reason == "CrashLoopBackOff",
		}
	}

	summary := &PodSummary{
		Name:        pod.Name,
		Namespace:   pod.Namespace,
		Phase:       pod.Status.Phase,
		Conditions:  pod.Status.Conditions,
		Annotations: pod.Annotations,
		Labels:      pod.Labels,
		Containers:  containers,
	}

	return summary
}
