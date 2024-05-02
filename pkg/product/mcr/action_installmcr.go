package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type installMCRStep struct {
	id string
}

func (s installMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-install", s.id)
}

func (s installMCRStep) Validate(_ context.Context) error {
	return nil
}

func (s installMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR install step", slog.String("ID", s.Id()))
	return nil
}
