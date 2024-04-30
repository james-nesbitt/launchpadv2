package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

// Install Controller, install worker, upgrade controller, upgrade worker

type activateK0sStep struct {
	id string
}

func (s activateK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-activate", s.id)
}

func (s activateK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s activate step", slog.String("ID", s.id))
	return nil
}
