package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type configureMCRStep struct {
	baseStep

	id string
}

func (s configureMCRStep) Id() string {
	return fmt.Sprintf("%s:mcr-configure", s.id)
}

func (s configureMCRStep) Validate(_ context.Context) error {
	return nil
}

func (s configureMCRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR configure step", slog.String("ID", s.id))
	return nil
}
