package mcr

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker/swarm"
)

// Requires declare that we need a HostsRoles dependency
func (p *MCR) Requires(_ context.Context) dependency.Requirements {
	p.rhr = host.NewHostsRolesRequirement(
		fmt.Sprintf("%s:%s", p.id, host.RequiresHostsRolesType),
		fmt.Sprintf("MCR '%s' needs hosts as installation targets, using roles: %s", p.id, strings.Join(MCRHostRoles, ",")),
		host.HostsRolesFilter{
			Roles: MCRHostRoles,
			Min:   1,
			Max:   0,
		},
	)

	return dependency.Requirements{p.rhr}
}

func (p *MCR) getHosts(ctx context.Context) (*host.Hosts, error) {
	if p.rhr == nil {
		return nil, fmt.Errorf("%s: Missing Hosts dependency", ComponentType)
	}

	d := p.rhr.Matched(ctx)
	if d == nil {
		return nil, fmt.Errorf("%s: Hosts dependency not Matched", ComponentType)
	}

	if err := d.Met(ctx); err != nil {
		return nil, fmt.Errorf("%s: Host dependency not Met: %w", ComponentType, err)
	}

	rhd, ok := d.(host.HostsDependency)
	if !ok {
		return nil, fmt.Errorf("%s: Hosts dependency wrong type: %+v", ComponentType, d)
	}

	hs := rhd.ProduceHosts(ctx)

	return hs, nil
}

// Provides dependencies
func (p *MCR) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if dhr, ok := r.(dockerhost.DockerHostsRequirement); ok {
		// DockerHosts dependency
		//
		// @TODO should we use a single dependency for this?
		//

		v := dhr.NeedsDockerHost(ctx)

		dhd := dockerhost.NewDockerHostsDependency(
			fmt.Sprintf("%s:%s", ComponentType, dockerhost.ImplementationType),
			fmt.Sprintf("Docker hosts for requirement: %s", r.Describe()),
			func(ctx context.Context) (*dockerhost.DockerHosts, error) {
				dh, err := p.getDockerHosts(ctx, v)
				if err != nil {
					return nil, err
				}

				return dh, nil
			},
		)

		return dhd, nil
	}
	if dsr, ok := r.(swarm.DockerSwarmRequirement); ok {
		// DockerSwarm dependency

		v := dsr.NeedsDockerSwarm(ctx)

		dsd := swarm.NewDockerSwarmDependency(
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

		return dsd, nil
	}

	return nil, nil
}

func (p *MCR) getDockerHosts(ctx context.Context, version docker.Version) (*dockerhost.DockerHosts, error) {
	hs, err := p.getHosts(ctx)
	if err != nil {
		return nil, err
	}

	return dockerhost.NewDockerHosts(
		hs,
		version,
	), nil
}
