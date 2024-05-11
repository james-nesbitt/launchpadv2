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

func (s *uninstallMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR uninstall step", slog.String("ID", s.Id()))

	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		i, ierr := h.Docker(ctx).Info(ctx)
		if ierr != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR not installed", h.Id()), slog.Any("host", h))
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: disabling and uninstalling MCR", h.Id()), slog.Any("host", h), slog.Any("info", i))

		if err := disableMCRService(ctx, h); err != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR disable failed", h.Id()), slog.Any("host", h))
			return err
		}
		if err := uninstallMCR(ctx, h); err != nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR remove failed", h.Id()), slog.Any("host", h))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func disableMCRService(ctx context.Context, h *dockerhost.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: disabling MCR services: %s", h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", h))

	if err := h.ServiceDisable(ctx, MCRServices); err != nil {
		return err
	}

	return nil
}

func uninstallMCR(ctx context.Context, h *dockerhost.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: uninstall MCR: %s", h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", h))

	if err := h.RemovePackages(ctx, MCRUninstallPackages); err != nil {
		return err
	}

	return nil
}
