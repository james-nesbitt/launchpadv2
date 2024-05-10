package mcr

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker/swarm"
)

// Requires declare that we need a HostsRoles dependency.
func (c *MCR) RequiresDependencies(_ context.Context) dependency.Requirements {
	if c.mhr == nil {
		c.mhr = host.NewHostsRolesRequirement(
			fmt.Sprintf("%s:%s:manager", host.HostsComponentType, c.Name()),
			fmt.Sprintf("MCR '%s' needs at least one manager host as installation targets, using roles: %s", c.id, strings.Join(MCRManagerHostRoles, ",")),
			host.HostsRolesFilter{
				Roles: MCRManagerHostRoles,
				Min:   1,
				Max:   0,
			},
		)
	}
	if c.whr == nil {
		c.whr = host.NewHostsRolesRequirement(
			fmt.Sprintf("%s:%s:worker", host.HostsComponentType, c.Name()),
			fmt.Sprintf("MCR '%s' accepts any number of worker hosts as installation targets, using roles: %s", c.id, strings.Join(MCRWorkerHostRoles, ",")),
			host.HostsRolesFilter{
				Roles: MCRWorkerHostRoles,
				Min:   0,
				Max:   0,
			},
		)
	}

	return dependency.Requirements{
		c.mhr,
		c.whr,
	}
}

// Provides dependencies.
func (c *MCR) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
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
				func(ctx context.Context) (dockerhost.Hosts, error) {
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
	if dsr, ok := r.(swarm.DockerSwarmRequirement); ok {
		// DockerSwarm dependency

		v := dsr.NeedsDockerSwarm(ctx)

		if c.dsd == nil {
			c.dsd = swarm.NewDockerSwarmDependency(
				fmt.Sprintf("%s:%s", ComponentType, "DockerSwarm"),
				fmt.Sprintf("Docker Swarm for requirement: %s", r.Describe()),
				func(ctx context.Context) (*swarm.Swarm, error) {
					cfg := swarm.Config{}
					s := swarm.NewSwarm(cfg)

					if err := s.ValidateVersion(v); err != nil {
						return nil, err
					}

					return s, nil
				},
			)
		}

		return c.dsd, nil
	}

	return nil, nil
}

// ValidateDockerVersion to say that this component can provided the requested version.
func (c MCR) ValidateDockerVersion(v dockerimplementation.Version) error {
	return nil
}
