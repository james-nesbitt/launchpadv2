package kubernetes

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

// ValidateK8sDependencyConfig validate a Kubernetes client request configuration.
func (c Component) ValidateK8sDependencyConfig(kc kubernetes.Version) error {
	return nil
}

func (c *Component) kubernetesImplementation(ctx context.Context) (*kubernetes.Kubernetes, error) {
	return kubernetes.NewKubernetes(c.config), nil
}
