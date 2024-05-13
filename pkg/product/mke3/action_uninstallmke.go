package mke3

import (
	"context"
	"fmt"
	"log/slog"

	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

type uninstallMKEStep struct {
	id string
}

func (s uninstallMKEStep) Id() string {
	return fmt.Sprintf("%s:mke3-uninstall", s.id)
}

func (s uninstallMKEStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE3 uninstall step", slog.String("ID", s.Id()))
	return nil
}

func mkeUninstall(ctx context.Context, h *dockerhost.Host, c Config) error {
	i, ierr := h.Docker(ctx).Info(ctx)
	if ierr != nil {
		return fmt.Errorf("%s: has no docker info, so can't be used to uninstall", h.Id())
	}
	if i.Swarm.Cluster == nil {
		return fmt.Errorf("%s: has no docker swarm cluster info, so can't be used to uninstall", h.Id())
	}

	cmd := []string{"run", "--volume='/var/run/docker.sock:/var/run/docker.sock'", c.bootstrapperImage(), "uninstall-ucp"}
	cmd = append(cmd, fmt.Sprintf("--id=%s", i.Swarm.Cluster.ID))
	cmd = append(cmd, c.bootStrapperUninstallArgs()...)

	o, e, err := h.Docker(ctx).Run(ctx, cmd, dockerimplementation.RunOptions{ShowOutput: true, ShowError: true})
	if err != nil {
		return fmt.Errorf("MKE bootstrap failed: %s :: %s", err.Error(), e)
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: MKE bootstrapper uninstall succeeded: %s", h.Id(), o))

	return nil
}

func mkePruneAfterInstall(ctx context.Context, hs dockerhost.Hosts) error {
	if err := hs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		o, e, err := h.Docker(ctx).Run(ctx, []string{"system", "prune", "--all", "--force"}, dockerimplementation.RunOptions{ShowOutput: true, ShowError: true})
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
