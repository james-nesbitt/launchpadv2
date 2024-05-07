package mcr

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Mirantis/launchpad/pkg/host"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

var (
	MCRInstallerPath = "/tmp/mcr-install.sh"
)

type installMCRStep struct {
	baseStep

	id string
}

func (s installMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-install", s.id)
}

func (s installMCRStep) Validate(_ context.Context) error {
	return nil
}

func (s installMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR install step", slog.String("ID", s.Id()))

	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	err := hs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		if err := s.downloadMCRInstaller(ctx, h); err != nil {
			return err
		}

		if err := s.runMCRInstaller(ctx, h); err != nil {
			return err
		}

		if err := s.enableMCRService(ctx, h); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s installMCRStep) downloadMCRInstaller(ctx context.Context, h *dockerhost.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: downloading MCR Installer: %s", h.Id(), s.c.config.InstallURLLinux), slog.Any("host", h))
	irs, igerr := http.Get(s.c.config.InstallURLLinux)
	if igerr != nil {
		return igerr
	}

	defer irs.Body.Close()

	if err := h.Upload(ctx, irs.Body, MCRInstallerPath, 0777, host.ExecOptions{Sudo: true}); err != nil {
		return err
	}

	return nil
}

func (s installMCRStep) runMCRInstaller(ctx context.Context, h *dockerhost.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: running MCR Installer: %s", h.Id(), s.c.config.InstallURLLinux), slog.Any("host", h))

	cmd := fmt.Sprintf("DOCKER_URL=%s CHANNEL=%s VERSION=%s bash %s", s.c.config.RepoURL, s.c.config.Channel, s.c.config.Version, MCRInstallerPath)
	if _, e, err := h.Exec(ctx, cmd, nil, host.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: mcr installer fail: %s \n %s", h.Id(), err.Error(), e)
	}
	return nil
}

func (s installMCRStep) enableMCRService(ctx context.Context, h *dockerhost.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: enabling MCR service: %s", h.Id(), s.c.config.InstallURLLinux), slog.Any("host", h))

	if err := h.ServiceEnable(ctx, []string{"docker"}); err != nil {
		return err
	}

	return nil
}
