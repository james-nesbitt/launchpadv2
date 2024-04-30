package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

type prepaseHostStep struct {
	id string
}

func (s prepaseHostStep) Id() string {
	return fmt.Sprintf("%s:k0s-prepare-hosts", s.id)
}

func (s prepaseHostStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s prepare-hosts step", slog.String("ID", s.id))
	return nil
}
