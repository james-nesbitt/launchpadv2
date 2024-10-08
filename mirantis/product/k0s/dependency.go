package k0s

import (
	"context"
	"fmt"

	kubeimpl "github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

// --- Dependency requirements ---

// Requires declare requirements needed for this component
//  1. at least one host with role "controller"
//  2. any workers with role "worker" (optional)
func (c *Component) RequiresDependencies(_ context.Context) dependency.Requirements {
	if c.hs == nil {

		c.hs = host.NewHostsFilterRequirement(
			fmt.Sprintf("%s:%s", host.ComponentType, c.Name()),
			fmt.Sprintf("K0S '%s' takes all nodes marked with K0S; needs at least one controller host as installation targets", c.Name()),
			func(ctx context.Context, hs host.Hosts) (host.Hosts, error) {
				// filter for any nodes with an K0S plugin
				// - we also check to make sure that there is at least one controller
				fhs := host.NewHosts()

				chc := 0 // count how many controllers we find
				for _, h := range hs {
					hk := HostGetK0s(h) // get the host K0s plugin
					if hk == nil {
						// node has no k0s plugin, so it is not a target for install
						continue
					}

					if hk.IsController() {
						chc = chc + 1
					}

					fhs.Add(h)
				}

				if chc == 0 {
					return fhs, fmt.Errorf("%s: no controllers in cluster", c.Name())
				}

				return fhs, nil
			},
		)

	}

	return dependency.Requirements{
		c.hs,
	}
}

// --- Dependency provides ---

// Provides dependencies
//  1. provide any Kubernetes dependencies if the versions match
//  2. provide any k0S dependencies if the version match
func (p *Component) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
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

	return nil, nil
}
