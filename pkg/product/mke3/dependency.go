package mke3

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/product/mke3/implementation"
)

// Requires declare that we need a HostsRoles dependency
func (c *MKE3) Requires(_ context.Context) dependency.Requirements {
	if c.dhr == nil {
		c.dhr = dockerhost.NewDockerHostsRequirement(
			fmt.Sprintf("%s:%s", c.Name(), dockerhost.ImplementationType),
			fmt.Sprintf("%s: Needs docker hosts to install on", ComponentType),
			docker.Version{},
		)
	}

	return dependency.Requirements{
		c.dhr,
	}
}

// Provides dependencies
func (c *MKE3) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if mke3r, ok := r.(implementation.MKE3Requirement); ok {
		// MKE3 dependency
		dc := mke3r.RequiresMKE3(ctx)

		if err := c.ValidateMKE3DependencyConfig(dc); err != nil {
			return nil, err
		}

		if c.mked == nil {
			c.mked = implementation.NewMKE3Dependency(
				fmt.Sprintf("%s:%s", ComponentType, implementation.ImplementationType),
				fmt.Sprintf("%s: MKE implementation from '%s'", ComponentType, c.id),
				func(_ context.Context) (*implementation.API, error) {
					return &implementation.API{}, nil
				},
			)
		}

		return c.mked, nil
	}
	if k8sr, ok := r.(kubernetes.KubernetesRequirement); ok {
		// Kubernetes dependency
		dc := k8sr.RequiresKubernetes(ctx)

		if err := c.ValidateK8sDependencyConfig(dc); err != nil {
			return nil, err
		}

		if c.k8sd == nil {
			c.k8sd = kubernetes.NewKubernetesDependency(
				fmt.Sprintf("%s:%s", ComponentType, kubernetes.ImplementationType),
				fmt.Sprintf("%s: kubernetes implementation from %s", ComponentType, c.id),
				func(ctx context.Context) (*kubernetes.Kubernetes, error) {
					return c.kubernetesImplementation(ctx)
				},
			)
		}

		return c.k8sd, nil
	}

	return nil, nil
}

// ValidateDependencyConfig validate that the dependency configuration can be met
func (c MKE3) ValidateMKE3DependencyConfig(dc implementation.MKE3DependencyConfig) error {
	return nil
}
