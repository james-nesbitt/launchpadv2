package kubernetes

import (
	"context"
	"fmt"

	kubeimpl "github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

// --- Dependency provides ---

// Provides dependencies
//  1. provide any Kubernetes dependencies if the versions match
func (c *Component) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if k8sr, ok := r.(kubeimpl.KubernetesRequirement); ok {
		// Kubernetes dependency

		kv := k8sr.RequiresKubernetes(ctx)

		if err := c.ValidateK8sDependencyConfig(kv); err != nil {
			// this kubernetes component doesn't meet the kubernetes needs of the requirement (version, apis etc)
			return nil, err
		}

		if c.k8sd == nil {
			c.k8sd = kubeimpl.NewKubernetesDependency(
				fmt.Sprintf("%s:%s", ComponentType, kubeimpl.ImplementationType),
				fmt.Sprintf("%s: kubernetes implementation: %s", ComponentType, r.Describe()),
				c.kubernetesImplementation, // use our kubernetes implementation constructor in the dependency
			)
		}

		return c.k8sd, nil
	}

	return nil, nil
}
