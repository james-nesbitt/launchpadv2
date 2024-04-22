package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type restartMCRStep struct {
	id string
}

func (s restartMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-restart", s.id)
}

func (s restartMCRStep) Validate(_ context.Context) error {
	return nil
}

func (s restartMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR restart step", slog.String("ID", s.id))
	return nil
}
