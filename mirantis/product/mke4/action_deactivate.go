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

func (s deactivateStep) Id() string {
	return fmt.Sprintf("%s:mke4-deactivate", s.id)
}

func (s deactivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 deactivate step", slog.String("ID", s.Id()))
	return nil
}
