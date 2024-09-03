package mcr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

type restartMCRStep struct {
	baseStep

	id string
}

func (s restartMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-restart", s.id)
}

func (s *restartMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR restart step", slog.String("ID", s.Id()))
	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		hm := HostGetMCR(h)
		return hm.RestartMCRService(ctx)
	}); err != nil {
		return err
	}

	return nil
}
