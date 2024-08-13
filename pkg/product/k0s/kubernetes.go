package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

// ValidateK8sDependencyConfig validate a Kubernetes client request configuration.
func (p Component) ValidateK8sDependencyConfig(kc kubernetes.Version) error {
	return nil
}

func (p *Component) kubernetesImplementation(ctx context.Context) (*kubernetes.Kubernetes, error) {
	kcc, kccerr := p.K0SKubeConfigAdmin(ctx)
	if kccerr != nil {
		return nil, kccerr
	}

	c := kubernetes.Config{
		KubeClientConfig: *kcc,
	}
	k := kubernetes.NewKubernetes(c)
	return k, nil
}
