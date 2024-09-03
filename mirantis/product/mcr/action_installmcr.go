package mcr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
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
		hm := HostGetMCR(h)

		i, ierr := hm.DockerInfo(ctx)

		if ierr == nil && i.ServerVersion == s.c.config.Version {
			slog.InfoContext(ctx, fmt.Sprintf("%s: MCR already at version %s", h.Id(), i.ServerVersion), slog.Any("host", h))
		} else {
			if i.ServerVersion == "" {
				slog.InfoContext(ctx, fmt.Sprintf("%s: installing MCR version %s", h.Id(), i.ServerVersion), slog.Any("host", h))
			} else {
				slog.InfoContext(ctx, fmt.Sprintf("%s: upgrading MCR from version %s -> %s", h.Id(), i.ServerVersion, s.c.config.Version), slog.Any("host", h))
			}

			if err := hm.DownloadMCRInstaller(ctx, s.c.config); err != nil {
				return err
			}

			if err := hm.RunMCRInstaller(ctx, s.c.config); err != nil {
				return err
			}

			if i.ServerVersion == "" {
				if err := hm.EnableMCRService(ctx); err != nil {
					return err
				}
			}

			if _, err := hm.DockerInfo(ctx); err != nil {
				return fmt.Errorf("%s: MCR discovery error after install", h.Id())
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
