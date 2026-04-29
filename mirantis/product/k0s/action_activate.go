package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/k0sproject/version"
)

/**
 * Activate K0S
 *
 * K0S must already be installed and configured. Here
 * we ensure that a cluster is up and running.
 *
 * If a host is already running k0s at the desired version, it is skipped.
 */

// hostNeedsInstall returns true if the host needs k0s installed or upgraded.
// Returns false only if k0s is already running at exactly the desired version.
func HostNeedsInstall(info HostDiscovery, desired version.Version) bool {
	if !info.Running {
		return true
	}
	if info.RunningVersion == nil {
		return true
	}
	return !info.RunningVersion.Equal(&desired)
}

type activateK0sStep struct {
	baseStep
	id string
}

func (s activateK0sStep) ID() string {
	return fmt.Sprintf("%s:k0s-activate", s.id)
}

func (s activateK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s activate step", slog.String("ID", s.id))

	slog.DebugContext(ctx, "Looking for leader in controllers")
	l := s.c.GetLeaderHost(ctx)
	if l == nil {
		return fmt.Errorf("could not find a leader")
	}

	baseCfg := s.c.config.K0sConfig
	csans := s.c.CollectClusterSans(ctx)

	lkh := HostGetK0s(l)

	if !HostNeedsInstall(s.c.state.HostInfo[l.ID()], s.c.config.Version) {
		slog.InfoContext(ctx, fmt.Sprintf("%s: leader k0s already at desired version, skipping", l.ID()))
	} else {
		// If k0s is running (upgrade case), stop it before reinstalling.
		// JoinCluster handles stop internally; we handle the leader manually here.
		if _, sterr := lkh.Status(ctx); sterr == nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: leader k0s running, stopping for install/upgrade", l.ID()))
			if err := lkh.K0sStop(ctx); err != nil {
				// Stop failure is non-fatal — k0s may be partially running.
				slog.WarnContext(ctx, fmt.Sprintf("%s: stop failed (continuing with install): %s", l.ID(), err.Error()))
			}
		} else {
			slog.DebugContext(ctx, fmt.Sprintf("%s: using as leader in new cluster", l.ID()))
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: writing config to leader host", l.ID()))
		if werr := lkh.BuildAndWriteK0sConfig(ctx, baseCfg, csans); werr != nil {
			return werr
		}

		if err := lkh.InstallNewCluster(ctx, s.c.config); err != nil {
			return fmt.Errorf("failed to install new cluster on leader %s: %w", l.ID(), err)
		}
	}

	chs, cherr := s.c.GetControllerHosts(ctx)
	if cherr != nil {
		return fmt.Errorf("could not retrieve controller hosts: %s", cherr.Error())
	}
	slog.InfoContext(ctx, "Sequentially adding controller hosts")
	if err := chs.Sequential(ctx, func(ctx context.Context, h *host.Host) error {
		if !HostNeedsInstall(s.c.state.HostInfo[h.ID()], s.c.config.Version) {
			slog.InfoContext(ctx, fmt.Sprintf("%s: k0s already at desired version, skipping", h.ID()))
			return nil
		}

		kh := HostGetK0s(h)

		slog.InfoContext(ctx, fmt.Sprintf("%s: writing config to controller host", h.ID()))
		if werr := kh.BuildAndWriteK0sConfig(ctx, baseCfg, csans); werr != nil {
			return werr
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: joining as controller to '%s' cluster", h.ID(), l.ID()))
		return kh.JoinCluster(ctx, l, RoleController, s.c.config)
	}, true); err != nil {
		return fmt.Errorf("error joining controller hosts: %s", err.Error())
	}

	whs, wherr := s.c.GetWorkerHosts(ctx)
	if wherr != nil {
		return fmt.Errorf("could not retrieve worker hosts: %s", wherr.Error())
	}
	slog.InfoContext(ctx, "In parallel adding worker hosts")
	if err := whs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		if !HostNeedsInstall(s.c.state.HostInfo[h.ID()], s.c.config.Version) {
			slog.InfoContext(ctx, fmt.Sprintf("%s: k0s already at desired version, skipping", h.ID()))
			return nil
		}

		eh := exec.HostGetExecutor(h)
		eh.Connect(ctx)
		kh := HostGetK0s(h)

		slog.InfoContext(ctx, fmt.Sprintf("%s: joining as worker to '%s' cluster", h.ID(), l.ID()))
		return kh.JoinCluster(ctx, l, RoleWorker, s.c.config)
	}); err != nil {
		return fmt.Errorf("error joining worker hosts: %s", err.Error())
	}

	return nil
}
