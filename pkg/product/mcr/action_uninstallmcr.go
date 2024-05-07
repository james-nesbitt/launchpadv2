package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallMCRStep struct {
	baseStep

	id string
}

func (s uninstallMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-uninstall", s.id)
}

func (s uninstallMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR uninstall step", slog.String("ID", s.Id()))
	return nil
}
