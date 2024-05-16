package mcr

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

var (
	MCRInstallerPath = "mcr-install"
	MCRServices      = []string{"docker", "containerd"}
)

type installMCRStep struct {
	baseStep

	id string
}

func (s installMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-install", s.id)
}

func (s *installMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR install step", slog.String("ID", s.Id()))

	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		i, ierr := getHostDocker(ctx, h).Info(ctx)

		if ierr == nil && i.ServerVersion == s.c.config.Version {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR already at version %s", h.Id(), i.ServerVersion), slog.Any("host", h))
		} else {
			if i.ServerVersion == "" {
				slog.InfoContext(ctx, fmt.Sprintf("%s: installing MCR version %s", h.Id(), i.ServerVersion), slog.Any("host", h))
			} else {
				slog.InfoContext(ctx, fmt.Sprintf("%s: upgrading MCR from version %s -> %s", h.Id(), i.ServerVersion, s.c.config.Version), slog.Any("host", h))
			}

			if err := s.downloadMCRInstaller(ctx, h); err != nil {
				return err
			}

			if err := s.runMCRInstaller(ctx, h); err != nil {
				return err
			}

			if i.ServerVersion == "" {
				if err := s.enableMCRService(ctx, h); err != nil {
					return err
				}
			}

			if _, err := getHostDocker(ctx, h).Info(ctx); err != nil {
				return fmt.Errorf("%s: MCR discovery error after install", h.Id())
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s installMCRStep) downloadMCRInstaller(ctx context.Context, h *host.Host) error {
	ir := s.c.config.InstallURLLinux

	slog.InfoContext(ctx, fmt.Sprintf("%s: downloading MCR Installer: %s", h.Id(), ir), slog.Any("host", h))
	irs, igerr := http.Get(ir)
	if igerr != nil {
		return igerr
	}

	defer irs.Body.Close()

	if err := exec.HostGetExecutor(h).Upload(ctx, irs.Body, MCRInstallerPath, 0777, exec.ExecOptions{}); err != nil {
		return err
	}

	return nil
}

func (s installMCRStep) runMCRInstaller(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: running MCR Installer (%s)", h.Id(), MCRInstallerPath), slog.Any("host", h))

	cmd := fmt.Sprintf("DOCKER_URL=%s CHANNEL=%s VERSION=%s bash %s", s.c.config.RepoURL, s.c.config.Channel, s.c.config.Version, MCRInstallerPath)
	if _, e, err := exec.HostGetExecutor(h).Exec(ctx, cmd, nil, exec.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: mcr installer fail: %s \n %s", h.Id(), err.Error(), e)
	}
	return nil
}

func (s installMCRStep) enableMCRService(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: enabling MCR services: %s", h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", h))

	if err := exec.HostGetExecutor(h).ServiceEnable(ctx, MCRServices); err != nil {
		return err
	}

	return nil
}
