package k0s

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	kubeimpl "github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/product/k0s/implementation"
)

// --- Dependency requirements ---

// Requires declare requirements needed for this component
//  1. at least one host with role "controller"
//  2. any workers with role "worker" (optional)
func (p K0S) Requires(_ context.Context) dependency.Requirements {
	if p.chr == nil {
		p.chr = host.NewHostsRolesRequirement(
			p.Name(),
			fmt.Sprintf("MCR '%s' needs hosts as controller installation targets, using role: %s", p.id, ControllerHostRole),
			host.HostsRolesFilter{
				Roles: []string{ControllerHostRole},
				Min:   1,
				Max:   0,
			},
		)
	}

	if p.whr == nil {
		p.whr = host.NewHostsRolesRequirement(
			p.Name(),
			fmt.Sprintf("MCR '%s' needs hosts as worker installation targets, using role: %s", p.id, WorkerHostRole),
			host.HostsRolesFilter{
				Roles: []string{WorkerHostRole},
				Min:   0,
				Max:   0,
			},
		)
	}

	return dependency.Requirements{
		p.chr,
		p.whr,
	}
}

// --- Dependency provides ---

// Provides dependencies
//  1. provide any Kubernetes dependencies if the versions match
//  2. provide any k0S dependencies if the version match
func (p *K0S) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if k8sr, ok := r.(kubeimpl.KubernetesRequirement); ok {
		// Kubernetes dependency

		c := k8sr.RequiresKubernetes(ctx)

		if err := p.ValidateK8sDependencyConfig(c); err != nil {
			return nil, err
		}

		if p.k8sd == nil {
			p.k8sd = kubeimpl.NewKubernetesDependency(
				fmt.Sprintf("%s:%s", ComponentType, kubeimpl.ImplementationType),
				fmt.Sprintf("%s: kubernetes implementation: %s", ComponentType, r.Describe()),
				func(ctx context.Context) (*kubeimpl.Kubernetes, error) {
					return p.kubernetesImplementation(ctx)
				},
			)
		}

		return p.k8sd, nil
	}

	if k0sr, ok := r.(implementation.RequiresK0s); ok {
		// Kubernetes dependency

		c := k0sr.RequireK0s(ctx)

		if err := p.ValidateK0sDependencyConfig(c); err != nil {
			return nil, err
		}

		if p.k0sd == nil {
			p.k0sd = implementation.NewK0sDependency(
				fmt.Sprintf("%s:%s", ComponentType, implementation.ImplementatonType),
				fmt.Sprintf("%s: kubernetes implementation: %s", ComponentType, r.Describe()),
				func(ctx context.Context) (*implementation.API, error) {
					return p.k0sImplementation(ctx)
				},
			)
		}

		return p.k8sd, nil
	}

	return nil, nil
}
