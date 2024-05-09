package mcr

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

var (
	MCRUninstallPackages = []string{"docker-ee", "docker-ee-cli", "containerd"}
)

type uninstallMCRStep struct {
	baseStep

	id string
}

func (s uninstallMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-uninstall", s.id)
}

func (s uninstallMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR uninstall step", slog.String("ID", s.Id()))

	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		hs := s.c.state.hosts.GetOrCreate(h.Id())

		hsi := hs.getInfo()

		if hsi.ServerVersion == "" {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR not installed", h.Id()), slog.Any("host", h))
		}

		if err := s.disableMCRService(ctx, h); err != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR install failed", h.Id()), slog.Any("host", h))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := s.updateDockerState(ctx); err == nil {
		slog.WarnContext(ctx, "MCR discovery still discovers MCR after uninstallation", slog.Any("error", err.Error()))
		return err
	}

	return nil
}

func (s uninstallMCRStep) disableMCRService(ctx context.Context, h *dockerhost.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: disabling MCR services: %s", h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", h))

	if err := h.ServiceDisable(ctx, MCRServices); err != nil {
		return err
	}

	return nil
}
