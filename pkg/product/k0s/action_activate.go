package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

/**
 * Activate K0S
 *
 * K0S must already be installed and configured. Here
 * we ensure that a cluster is up and running.
 *
 * @TODO we don't yet do any version detection/comparison (desired vs running)
 */

type activateK0sStep struct {
	baseStep
	id string
}

func (s activateK0sStep) Id() string {
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

	if ls, lserr := lkh.Status(ctx); lserr == nil {
		slog.DebugContext(ctx, fmt.Sprintf("%s: discovered as leader", l.Id()), slog.Any("status", ls))
	} else {
		// leader has no k0s running, so start a new cluster
		slog.DebugContext(ctx, fmt.Sprintf("%s: using as leader in new cluster", l.Id()), slog.Any("status", ls))

		lkh := HostGetK0s(l)

		slog.InfoContext(ctx, fmt.Sprintf("%s: writing config to leader host", l.Id()))
		if werr := lkh.BuildAndWriteK0sConfig(ctx, baseCfg, csans); werr != nil {
			return werr
		}

		lkh.InstallNewCluster(ctx, s.c.config)
	}

	chs, cherr := s.c.GetControllerHosts(ctx)
	if cherr != nil {
		return fmt.Errorf("could not retrieve controller hosts: %s", cherr.Error())
	}
	if err := chs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		kh := HostGetK0s(h)

		slog.InfoContext(ctx, fmt.Sprintf("%s: writing config to controller host", h.Id()))
		if werr := kh.BuildAndWriteK0sConfig(ctx, baseCfg, csans); werr != nil {
			return werr
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: joining as controller to '%s' cluster", h.Id(), l.Id()))
		return kh.JoinCluster(ctx, l, RoleController, s.c.config)
	}); err != nil {
		return fmt.Errorf("error joining controller hosts: %s", err.Error())
	}

	whs, wherr := s.c.GetWorkerHosts(ctx)
	if wherr != nil {
		return fmt.Errorf("could not retrieve worker hosts: %s", wherr.Error())
	}
	if err := whs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		eh := exec.HostGetExecutor(h)
		eh.Connect(ctx)
		kh := HostGetK0s(h)

		slog.InfoContext(ctx, fmt.Sprintf("%s: joining as worker to '%s' cluster", h.Id(), l.Id()))
		return kh.JoinCluster(ctx, l, RoleWorker, s.c.config)
	}); err != nil {
		return fmt.Errorf("error joining worker hosts: %s", err.Error())
	}

	return nil
}
