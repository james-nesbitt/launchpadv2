package mcr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
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

func (s *uninstallMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR uninstall step", slog.String("ID", s.Id()))

	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		hm := HostGetMCR(h)

		i, ierr := hm.DockerInfo(ctx)
		if ierr != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR not installed", h.Id()), slog.Any("host", h))
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: disabling and uninstalling MCR", h.Id()), slog.Any("host", h), slog.Any("info", i))

		if err := hm.DisableMCRService(ctx); err != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR disable failed", h.Id()), slog.Any("host", h))
			return err
		}
		if err := hm.UninstallMCR(ctx); err != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR remove failed", h.Id()), slog.Any("host", h))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
