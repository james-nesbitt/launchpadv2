package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

type installK0sStep struct {
	id string
}

func (s installK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-install", s.id)
}

func (s installK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s install step", slog.String("ID", s.id))
	return nil
}
