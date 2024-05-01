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

// Requires declare that we need a HostsRoles dependency.
func (p *MCR) Requires(_ context.Context) dependency.Requirements {
	if p.mrhr == nil {
		p.mrhr = host.NewHostsRolesRequirement(
			fmt.Sprintf("%s:%s:manager", host.HostsComponentType, p.Name()),
			fmt.Sprintf("MCR '%s' needs at least one manager host as installation targets, using roles: %s", p.id, strings.Join(MCRManagerHostRoles, ",")),
			host.HostsRolesFilter{
				Roles: MCRManagerHostRoles,
				Min:   1,
				Max:   0,
			},
		)
	}
	if p.wrhr == nil {
		p.wrhr = host.NewHostsRolesRequirement(
			fmt.Sprintf("%s:%s:worker", host.HostsComponentType, p.Name()),
			fmt.Sprintf("MCR '%s' accepts any number of worker hosts as installation targets, using roles: %s", p.id, strings.Join(MCRWorkerHostRoles, ",")),
			host.HostsRolesFilter{
				Roles: MCRWorkerHostRoles,
				Min:   0,
				Max:   0,
			},
		)
	}

	return dependency.Requirements{
		p.mrhr,
		p.wrhr,
	}
}

// Provides dependencies.
func (p *MCR) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if dhr, ok := r.(dockerhost.DockerHostsRequirement); ok {
		// DockerHosts dependency
		//
		// @TODO should we use a single dependency for this?
		//

		v := dhr.NeedsDockerHost(ctx)

		if p.dhd == nil {
			p.dhd = dockerhost.NewDockerHostsDependency(
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
		}

		return p.dhd, nil
	}
	if dsr, ok := r.(swarm.DockerSwarmRequirement); ok {
		// DockerSwarm dependency

		v := dsr.NeedsDockerSwarm(ctx)

		if p.dsd == nil {
			p.dsd = swarm.NewDockerSwarmDependency(
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

		return p.dsd, nil
	}

	return nil, nil
}

func (p *MCR) getDockerHosts(ctx context.Context, version docker.Version) (*dockerhost.DockerHosts, error) {
	hs := host.Hosts{}

	addRequirementHosts := func(r dependency.Requirement) error {
		if r == nil {
			return fmt.Errorf("no requirement provided")
		}

		d := r.Matched(ctx)
		if d == nil {
			return fmt.Errorf("no hosts dependency matched for requirement")
		}

		hd, ok := d.(host.HostsDependency)
		if !ok {
			return fmt.Errorf("wrong dependency type")
		}

		for id, h := range *hd.ProduceHosts(ctx) {
			hs[id] = h
		}

		return nil
	}

	if err := addRequirementHosts(p.mrhr); err != nil {
		return nil, err
	}
	if err := addRequirementHosts(p.wrhr); err != nil {
		return nil, err
	}

	return dockerhost.NewDockerHosts(
		hs,
		version,
	), nil
}
