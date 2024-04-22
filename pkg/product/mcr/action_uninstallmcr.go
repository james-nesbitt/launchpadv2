package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallMCRStep struct {
	id string
}

func (s uninstallMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-uninstall", s.id)
}

func (s uninstallMCRStep) Validate(_ context.Context) error {
	return nil
}

func (s uninstallMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR uninstall step", slog.String("ID", s.id))
	return nil
}
