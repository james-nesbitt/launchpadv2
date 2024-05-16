package mcr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

type discoverStep struct {
	baseStep

	id          string
	failOnError bool
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:mcr-discover", s.id)
}

func (s *discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR discover step", slog.String("ID", s.Id()))

	hs, gherr := s.c.GetAllHosts(ctx)
	if gherr != nil {
		return fmt.Errorf("MCR has no hosts to discover: %s", gherr.Error())
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: discovering MCR state", h.Id()), slog.Any("host", h))

		if _, err := getHostDocker(ctx, h).Info(ctx); err != nil {
			slog.DebugContext(ctx, fmt.Sprintf("%s: MCR state discovery failure", h.Id()), slog.Any("host", h), slog.Any("error", err))
			return fmt.Errorf("%s: failed to update docker info: %s", h.Id(), err.Error())
		}

		return nil
	}); err != nil {
		if s.failOnError {
			return fmt.Errorf("docker info update failed: %s", err.Error())
		}
		slog.InfoContext(ctx, fmt.Sprintf("At least one host was missing MCR installation: %s", err.Error()))
	}

	slog.InfoContext(ctx, "MCR discovery succeeded")
	return nil
}
