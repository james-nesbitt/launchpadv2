package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:k0s-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s discover step", slog.String("ID", s.id))
	return nil
}
