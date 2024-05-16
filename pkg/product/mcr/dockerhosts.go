package mcr

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
)

// getHostDocker get a DockerExec instance for a host (nil if not possible, which should never happen)
func getHostDocker(h *host.Host) *dockerimplementation.DockerExec {
	e := exec.HostGetExecutor(h)
	if e == nil {
		return nil
	}
	//e.Connect(ctx)

	m := HostGetMCR(h)
	if m == nil {
		return nil
	}

	def := func(ctx context.Context, cmd string, i io.Reader, rops dockerimplementation.RunOptions) (string, string, error) {
		hopts := exec.ExecOptions{Sudo: m.SudoDocker()}

		if rops.ShowOutput {
			hopts.OutputLevel = "info"
		}
		if rops.ShowError {
			hopts.ErrorLevel = "warn"
		}

		return e.Exec(ctx, cmd, i, hopts)
	}

	return dockerimplementation.NewDockerExec(def)
}

// GetManagerHosts get the docker hosts for managers.
func (c MCR) GetManagerHosts(ctx context.Context) (host.Hosts, error) {
	hs, err := c.GetAllHosts(ctx)
	if err != nil {
		return hs, fmt.Errorf("MCR manager hosts retrieval error; %w", err)
	}

	mhs := host.NewHosts()
	for _, h := range hs {
		mcrh := HostGetMCR(h)

		if !mcrh.IsManager() {
			continue
		}

		mhs.Add(h)
	}

	return mhs, nil
}

// GetWorkerHosts get the docker hosts for workers.
func (c MCR) GetWorkerHosts(ctx context.Context) (host.Hosts, error) {
	hs, err := c.GetAllHosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("MCR worker hosts retrieval error; %w", err)
	}

	whs := host.NewHosts()
	for _, h := range hs {
		mh := HostGetMCR(h)

		if mh.IsManager() {
			continue
		}

		whs.Add(h)
	}

	return whs, nil
}

// GetAllHosts get the docker hosts for all hosts.
func (c MCR) GetAllHosts(ctx context.Context) (host.Hosts, error) {
	ghs, err := getRequirementHosts(ctx, c.hr)
	if err != nil {
		return nil, fmt.Errorf("hosts retrieval error; %w", err)
	}
	if len(ghs) == 0 {
		return nil, fmt.Errorf("MCR has no hosts to install on; %w", err)
	}

	hs := host.NewHosts()
	for _, h := range ghs {
		m := HostGetMCR(h)
		if m == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host provided to MCR has no MCR plugin", h.Id()))
			continue
		}

		if e := exec.HostGetExecutor(h); e == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host provided to MCR has no exec plugin", h.Id()))
			continue
		}

		if p := exec.HostGetPlatform(h); p == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host provided to MCR has no platform plugin", h.Id()))
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
//  4. get the hosts from the dependency and convert to DockerHosts
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
