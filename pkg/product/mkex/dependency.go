package mkex

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

// RequiresDependencies declare what dependencies we have
func (c *Component) RequiresDependencies(ctx context.Context) dependency.Requirements {
	if c.hsr == nil {

		c.hsr = host.NewHostsFilterRequirement(
			fmt.Sprintf("%s:%s:mcr", host.ComponentType, c.Name()),
			fmt.Sprintf("MCR '%s' takes all hosts marked as MKEX hosts", c.id),
			func(ctx context.Context, hs host.Hosts) (host.Hosts, error) {
				// filter for any nodes with an MCR plugin
				// - we also check to make sure that there is at least one manager

				fhs := host.NewHosts()

				if len(hs) == 0 {
					return fhs, fmt.Errorf("%s: no hosts in cluster", c.Name())
				}

				for _, h := range hs {
					mhp := HostGetMKEX(h)
					if mhp == nil {
						slog.WarnContext(ctx, fmt.Sprintf("%s: host '%s' is not an MKEX host, ignoring", c.Name(), h.Id()))
						continue
					}

					fhs.Add(h)
				}

				return fhs, nil
			},
		)

	}

	return dependency.Requirements{
		c.hsr,
	}
}

// Provides dependencies.
func (c *Component) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if dhr, ok := r.(dockerhost.DockerHostsRequirement); ok {
		// DockerHosts dependency
		//
		// @TODO should we use a single dependency for this?
		//

		v := dhr.NeedsDockerHost(ctx)
		if err := c.ValidateDockerVersion(v); err != nil {
			return nil, fmt.Errorf("%w; MCR '%s' doesn't provide the requested docker version: %+v", dependency.ErrRequirementNotMatched, c.Name(), err.Error())
		}

		if c.dhd == nil {
			c.dhd = dockerhost.NewDockerHostsDependency(
				fmt.Sprintf("%s:%s", ComponentType, dockerhost.ImplementationType),
				fmt.Sprintf("Docker hosts for requirement: %s", r.Describe()),
				func(ctx context.Context) (host.Hosts, error) {
					dh, err := c.GetAllHosts(ctx)
					if err != nil {
						return nil, err
					}

					return dh, nil
				},
			)
		}

		return c.dhd, nil
	}

	return nil, nil
}

// ValidateDockerVersion to say that this component can provided the requested version.
func (c Component) ValidateDockerVersion(v dockerimplementation.Version) error {
	return nil
}
