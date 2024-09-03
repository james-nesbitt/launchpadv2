package mke3

import (
	"context"
	"fmt"
	"log/slog"

	dockertypessystem "github.com/docker/docker/api/types/system"

	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/host"
)

type uninstallMKEStep struct {
	baseStep
	id string
}

func (s uninstallMKEStep) Id() string {
	return fmt.Sprintf("%s:mke3-uninstall", s.id)
}

func (s uninstallMKEStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "MKE3 Uninstall starting. Looking for a manager to operate on")

	mhs, gherr := s.c.GetManagerHosts(ctx)
	if gherr != nil {
		return fmt.Errorf("MCR has no hosts to discover: %s", gherr.Error())
	}

	var m *host.Host
	var mi dockertypessystem.Info

	for _, h := range mhs {
		i, ierr := dockerhost.HostGetDockerExec(h).Info(ctx)
		if ierr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host is not a docker machine", h.Id()), slog.Any("host", h))
			continue
		}

		if i.Swarm.ControlAvailable {
			slog.DebugContext(ctx, fmt.Sprintf("%s: this host can act as a leader.", h.Id()), slog.Any("host", h))

			m = h
			mi = i

			break
		}
		slog.DebugContext(ctx, fmt.Sprintf("%s: host rejected as leader", h.Id()), slog.Any("info", i))
	}

	if m == nil {
		return fmt.Errorf("no swarm leader found")
	}

	slog.DebugContext(ctx, fmt.Sprintf("%s: Leader found for un-installation", m.Id()), slog.Any("info", mi))

	if err := mkeUninstall(ctx, m, s.c.config); err != nil {
		return fmt.Errorf("%s: uninstall fail: %s", m.Id(), err.Error())
	}

	slog.InfoContext(ctx, "Pruning nodes after MKE uninstall")
	if err := mkePruneAfterInstall(ctx, mhs); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("system prune after uninstall failed: %s", err.Error()))
	}

	slog.InfoContext(ctx, "MKE3 Uninstall completed")

	return nil
}

func mkeUninstall(ctx context.Context, h *host.Host, c Config) error {
	d := dockerhost.HostGetDockerExec(h)
	if d == nil {
		return fmt.Errorf("%s: has has no docker plugin, so can't be used to uninstall", h.Id())
	}

	i, ierr := d.Info(ctx)
	if ierr != nil {
		return fmt.Errorf("%s: has no docker info, so can't be used to uninstall", h.Id())
	}
	if i.Swarm.Cluster == nil {
		return fmt.Errorf("%s: has no docker swarm cluster info, so can't be used to uninstall", h.Id())
	}

	cmd := []string{"run", "--volume='/var/run/docker.sock:/var/run/docker.sock'", c.bootstrapperImage(), "uninstall-ucp"}
	cmd = append(cmd, fmt.Sprintf("--id=%s", i.Swarm.Cluster.ID))
	cmd = append(cmd, c.bootStrapperUninstallArgs()...)

	o, e, err := d.Run(ctx, cmd, dockerimplementation.RunOptions{ShowOutput: true, ShowError: true})
	if err != nil {
		return fmt.Errorf("MKE bootstrap failed: %s :: %s", err.Error(), e)
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: MKE bootstrapper uninstall succeeded: %s", h.Id(), o))

	return nil
}

func mkePruneAfterInstall(ctx context.Context, hs host.Hosts) error {
	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		o, e, err := dockerhost.HostGetDockerExec(h).Run(ctx, []string{"system", "prune", "--all", "--force"}, dockerimplementation.RunOptions{ShowOutput: true, ShowError: true})
		if err != nil {
			return fmt.Errorf("%s: prune failed: %s :: %s", h.Id(), err.Error(), e)
		}

		slog.DebugContext(ctx, fmt.Sprintf("%s: prune: %s", h.Id(), o))

		return nil
	}); err != nil {
		return fmt.Errorf("MKE prune after uninstall failed: %s", err.Error())
	}

	return nil
}
