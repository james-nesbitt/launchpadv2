package mke3

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

// GetManagerHosts get the docker hosts for managers.
func (c MKE3) GetManagerHosts(ctx context.Context) (dockerhost.Hosts, error) {
	mhs := dockerhost.Hosts{}

	if hs, err := c.GetAllHosts(ctx); err != nil {
		return mhs, fmt.Errorf("could not retrieve managers: %s", err.Error())
	} else {
		for _, h := range hs {
			if h.IsSwarmManager(ctx) {
				mhs.Add(h)
			}
		}
	}

	return mhs, nil
}

// GetWorkerHosts get the docker hosts for workers.
func (c MKE3) GetWorkerHosts(ctx context.Context) (dockerhost.Hosts, error) {
	whs := dockerhost.Hosts{}

	if hs, err := c.GetAllHosts(ctx); err != nil {
		return whs, fmt.Errorf("could not retrieve managers: %s", err.Error())
	} else {
		for _, h := range hs {
			if !h.IsSwarmManager(ctx) {
				whs.Add(h)
			}
		}
	}

	return whs, nil
}

// GetAllHosts get the docker hosts for all hosts.
func (c MKE3) GetAllHosts(ctx context.Context) (dockerhost.Hosts, error) {
	return getRequirementDockerHosts(ctx, c.dhr)
}

// getRequirementDockerHosts retrieve the matching dockerhosts from the hosts requirement
//
// This needs to go through the following steps/checks:
//  1. is the requirement nil
//  2. was the requirement matched with a dependency
//  3. was the requirement matched with the right kind of dependency
//  4. get the hosts from the dependency and convert to DockerHosts
func getRequirementDockerHosts(ctx context.Context, r dependency.Requirement) (dockerhost.Hosts, error) {
	if r == nil {
		return nil, fmt.Errorf("requirement empty")
	}
	mhd := r.Matched(ctx)
	if mhd == nil {
		return nil, fmt.Errorf("requirement not matched")
	}
	mhddh, ok := mhd.(dockerhost.DockerHostsDependency) // check that we have a hosts dependency
	if !ok {
		// this should never happen, but it is possible
		return nil, fmt.Errorf("%w; %s Dependency is the wrong type", dependency.ErrDependencyNotMatched, mhd.Id())
	}

	hs := mhddh.ProvidesDockerHost(ctx) // get the Hosts

	return hs, nil
}
