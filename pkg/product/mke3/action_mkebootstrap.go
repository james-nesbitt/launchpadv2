package mke3

import (
	"context"
	"fmt"
	"log/slog"

	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

type runMKEBootstrapper struct {
	baseStep
	id string
}

func (s runMKEBootstrapper) Id() string {
	return fmt.Sprintf("%s:mke3-bootstrap", s.id)
}

func (s *runMKEBootstrapper) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE3 Bootstrapper step", slog.String("ID", s.Id()))

	hs, gerr := s.c.GetManagerHosts(ctx)
	if gerr != nil {
		return fmt.Errorf("MKE3 bootstrap step could not retrieve manager list: %s", gerr.Error())
	}
	if len(hs) == 0 {
		return fmt.Errorf("MKE3 bootstrap step received an empty manager list: %s", gerr.Error())
	}

	for _, h := range hs {
		err := mkeInstall(ctx, h, s.c.config)
		if err == nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MKE bootstrap succeeded", h.Id()))
			break
		}
		slog.ErrorContext(ctx, fmt.Sprintf("%s: install failed: %s", h.Id(), err.Error()))
	}

	return nil
}

func mkeInstall(ctx context.Context, h *dockerhost.Host, c Config) error {
	cmd := []string{"run", "--volume='/var/run/docker.sock:/var/run/docker.sock'", c.bootstrapperImage()}
	cmd = append(cmd, c.bootStrapperInstallArgs()...)

	o, e, err := h.Docker(ctx).Run(ctx, cmd, dockerimplementation.RunOptions{ShowOutput: true, ShowError: true})
	if err != nil {
		return fmt.Errorf("MKE bootstrap failed: %s :: %s", err.Error(), e)
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: MKE bootstrapper install succeeded: %s", h.Id(), o))

	return nil
}
