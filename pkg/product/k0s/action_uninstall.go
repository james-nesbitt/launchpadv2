package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallK0sStep struct {
	id string
}

func (s uninstallK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-uninstall", s.id)
}

func (s uninstallK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s uninstall step", slog.String("ID", s.id))
	return nil
}
