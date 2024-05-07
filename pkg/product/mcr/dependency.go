package mcr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

// GetManagerHosts get the docker hosts for managers.
func (c MCR) GetManagerHosts(ctx context.Context) (dockerhost.Hosts, error) {
	hs, err := getRequirementHosts(ctx, c.mhr)
	if err != nil {
		return nil, fmt.Errorf("Manager hosts retrieval error; %w", err)
	}
	if len(hs) == 0 {
		return nil, fmt.Errorf("MCR manager dependency has no hosts")
	}

	return hs, nil
}

// GetWorkerHosts get the docker hosts for workers.
func (c MCR) GetWorkerHosts(ctx context.Context) (dockerhost.Hosts, error) {
	hs, err := getRequirementHosts(ctx, c.whr)
	if err != nil {
		return nil, fmt.Errorf("Worker hosts retrieval error; %w", err)
	}
	if len(hs) == 0 {
		slog.WarnContext(ctx, "MCR worker dependency has no hosts.")
	}

	return hs, nil
}

// GetAllHosts get the docker hosts for all hosts.
func (c MCR) GetAllHosts(ctx context.Context) (dockerhost.Hosts, error) {
	errs := []error{}
	hs := dockerhost.Hosts{}

	if mhs, err := c.GetManagerHosts(ctx); err != nil {
		errs = append(errs, err)
	} else {
		hs.Merge(mhs)
	}
	if whs, err := c.GetWorkerHosts(ctx); err != nil {
		errs = append(errs, err)
	} else {
		hs.Merge(whs)
	}

	if len(errs) > 0 {
		return hs, errors.Join(errs...)
	}
	return hs, nil
}

// getRequirementHosts retrieve the matching docker hosts f the hosts requirement
//
// This needs to go through the following steps/checks:
//  1. is the requirement nil
//  2. was the requirement matched with a dependency
//  3. was the requirement matched with the right kind of dependency
//  4. get the hosts from the dependency and convert to DockerHosts
func getRequirementHosts(ctx context.Context, r dependency.Requirement) (dockerhost.Hosts, error) {
	if r == nil {
		return nil, fmt.Errorf("requirement empty")
	}
	mhd := r.Matched(ctx)
	if mhd == nil {
		return nil, fmt.Errorf("requirement not matched")
	}
	mhddh, ok := mhd.(host.HostsDependency) // check that we have a hosts dependency
	if !ok {
		// this should never happen, but it is possible
		return nil, fmt.Errorf("%w; %s Dependency is the wrong type", dependency.ErrDependencyNotMatched, mhd.Id())
	}

	hs := mhddh.ProduceHosts(ctx) // get the Hosts

	return dockerhost.NewDockerHosts(hs, dockerhost.HostsOptions{SudoDocker: true}), nil // convert the hosts to Docker hosts
}
