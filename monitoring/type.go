package monitoring

import v1 "k8s.io/api/core/v1"

type NodeSummary struct {
	Name        string
	ReadyStatus v1.ConditionStatus
	Conditions  []v1.NodeCondition
	Annotations map[string]string
	Labels      map[string]string
	Taints      []v1.Taint
	Capacity    v1.ResourceList
	Allocatable v1.ResourceList
}

type PodSummary struct {
	Name        string
	Namespace   string
	Phase       v1.PodPhase
	Conditions  []v1.PodCondition
	Annotations map[string]string
	Labels      map[string]string
	Containers  []ContainerSummary
}

type ContainerSummary struct {
	Name           string
	Usage          v1.ResourceList
	IsCrashLooping bool
}
