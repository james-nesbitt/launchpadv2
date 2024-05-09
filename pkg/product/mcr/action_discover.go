package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	baseStep

	id          string
	failOnError bool
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:mcr-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR discover step", slog.String("ID", s.Id()))

	if err := s.updateDockerState(ctx); err != nil {
		if s.failOnError {
			slog.WarnContext(ctx, "MCR discovery failed to discover installation", slog.Any("error", err.Error()))
			return err
		}
	}
	return nil
}
