package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

// GetManagerHosts get the docker hosts for managers.
func (c Component) GetManagerHosts(ctx context.Context) (host.Hosts, error) {
	hs, err := c.GetAllHosts(ctx)
	if err != nil {
		return hs, fmt.Errorf("K0S manager hosts retrieval error; %w", err)
	}

	mhs := host.NewHosts()
	for _, h := range hs {
		k0sh := HostGetK0S(h)
		if !k0sh.IsController() {
			continue
		}

		mhs.Add(h)
	}

	return mhs, nil
}

// GetWorkerHosts get the docker hosts for workers.
func (c Component) GetWorkerHosts(ctx context.Context) (host.Hosts, error) {
	hs, err := c.GetAllHosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("K0S worker hosts retrieval error; %w", err)
	}

	whs := host.NewHosts()
	for _, h := range hs {
		mh := HostGetK0S(h)
		if mh.IsController() {
			continue
		}

		whs.Add(h)
	}

	return whs, nil
}

// GetAllHosts get the docker hosts for all hosts.
func (c Component) GetAllHosts(ctx context.Context) (host.Hosts, error) {
	ghs, err := getRequirementHosts(ctx, c.hs)
	if err != nil {
		return nil, fmt.Errorf("hosts retrieval error; %w", err)
	}
	if len(ghs) == 0 {
		return nil, fmt.Errorf("K0S has no hosts to install on; %w", err)
	}

	hs := host.NewHosts()
	for _, h := range ghs {
		if m := HostGetK0S(h); m == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host provided to K0S has no K0S plugin", h.Id()))
			continue
		}

		if e := exec.HostGetExecutor(h); e == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host provided to K0S has no exec plugin", h.Id()))
			continue
		}

		if p := exec.HostGetPlatform(h); p == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host provided to K0S has no platform plugin", h.Id()))
			continue
		}

		hs.Add(h)
	}

	return hs, nil
}

// getRequirementHosts retrieve the matching docker hosts f the hosts requirement
//
// This needs to go through the following steps/checks:
//  1. is the requirement nil
//  2. was the requirement matched with a dependency
//  3. was the requirement matched with the right kind of dependency
func getRequirementHosts(ctx context.Context, r dependency.Requirement) (host.Hosts, error) {
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

	return mhddh.ProduceHosts(ctx) // get the Hosts
}
