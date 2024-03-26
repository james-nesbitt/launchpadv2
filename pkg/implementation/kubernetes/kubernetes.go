package kubernetes

import "context"

// NewKubernetes constructor
func NewKubernetes() Kubernetes {
	return Kubernetes{}
}

// Kubernetes implementation
type Kubernetes struct {
}

// --- Kubernetes dependencies ---

type ImplementKubernetes func(context.Context) (Kubernetes, error)
type ReceiveKubernetesImplementation func(context.Context, Kubernetes) error

// KubernetesRequires dependency requirement for Kubernetes
type KubernetesRequirement struct {
	Version  string // version constraint string
	Receiver ReceiveKubernetesImplementation
}

// KubernetesDependency depndency on kubernetes implementation
type KubernetesDependency struct {
	Description string
	Factory     ImplementKubernetes
}
