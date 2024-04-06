package k0s

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

// --- Dependency requirements ---

// Name for the component
func (p K0S) Requires(_ context.Context) dependency.Requirements {
	return dependency.Requirements{}
}

// --- Dependency provides ---

// Provides dependencies
func (p *K0S) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if k8sr, ok := r.(kubernetes.KubernetesRequirement); ok {
		// Kubernetes dependency

		c := k8sr.RequiresKubernetes(ctx)

		if err := p.ValidateK8sDependencyConfig(c); err != nil {
			return nil, err
		}

		d := kubernetes.NewKubernetesDependency(
			fmt.Sprintf("%s:%s", ComponentType, kubernetes.ImplementationType),
			fmt.Sprintf("%s: kubernetes implementation: %s", ComponentType, r.Describe()),
			func(ctx context.Context) (*kubernetes.Kubernetes, error) {
				return p.kubernetesImplementation(ctx)
			},
		)

		return d, nil
	}

	return nil, nil
}
