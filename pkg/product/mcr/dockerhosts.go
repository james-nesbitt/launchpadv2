package mcr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

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
