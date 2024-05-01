package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

// ValidateK8sDependencyConfig validate a Kubernetes client request configuration.
func (p K0S) ValidateK8sDependencyConfig(kc kubernetes.Version) error {
	return nil
}

func (p *K0S) kubernetesImplementation(_ context.Context) (*kubernetes.Kubernetes, error) {
	c := kubernetes.Config{}
	k := kubernetes.NewKubernetes(c)
	return k, nil
}
