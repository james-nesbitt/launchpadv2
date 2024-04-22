package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

type configureK0sStep struct {
	id string
}

func (s configureK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-configure", s.id)
}

func (s configureK0sStep) Validate(_ context.Context) error {
	return nil
}

func (s configureK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s configure step", slog.String("ID", s.id))
	return nil
}
