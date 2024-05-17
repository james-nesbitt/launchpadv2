package mcr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
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

	if err := hs.Each(ctx, s.restartMCRService); err != nil {
		return err
	}

	return nil
}

func (s restartMCRStep) restartMCRService(ctx context.Context, h *host.Host) error {
	if err := exec.HostGetExecutor(h).ServiceRestart(ctx, []string{"docker"}); err != nil {
		return err
	}

	return nil
}
