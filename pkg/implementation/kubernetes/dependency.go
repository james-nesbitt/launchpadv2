package kubernetes

import (
	"context"
)

// --- Kubernetes dependencies ---

type ImplementKubernetes func(context.Context) (Kubernetes, error)
type ReceiveKubernetesImplementation func(context.Context, Kubernetes) error

// KubernetesRequires dependency requirement for Kubernetes
type KubernetesRequirement struct {
	Version  string // version constraint string
	Receiver ReceiveKubernetesImplementation
}
