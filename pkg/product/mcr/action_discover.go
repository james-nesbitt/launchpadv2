package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:mcr-discover", s.id)
}

func (s discoverStep) Validate(_ context.Context) error {
	return nil
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR discover step", slog.String("ID", s.Id()))
	return nil
}
