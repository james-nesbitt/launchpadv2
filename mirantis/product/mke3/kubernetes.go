package mke3

import (
	"context"

	"github.com/Mirantis/launchpad/implementation/kubernetes"
)

// ValidateK8sDependencyConfig validate a Kubernetes client request configuration.
func (c Component) ValidateK8sDependencyConfig(kc kubernetes.Version) error {
	return nil
}

func (c *Component) kubernetesImplementation(_ context.Context) (*kubernetes.Kubernetes, error) {
	cf := kubernetes.Config{}
	k := kubernetes.NewKubernetes(cf)
	return k, nil
}
