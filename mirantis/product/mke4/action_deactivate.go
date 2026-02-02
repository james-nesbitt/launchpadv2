package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type deactivateStep struct {
	baseStep
	id string
}

func (s deactivateStep) ID() string {
	return fmt.Sprintf("%s:mke4-deactivate", s.id)
}

func (s deactivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 deactivate step", slog.String("ID", s.ID()))
	return nil
}
