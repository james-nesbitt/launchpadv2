package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type activateStep struct {
	baseStep
	id string
}

func (s activateStep) Id() string {
	return fmt.Sprintf("%s:activate-mke4", s.id)
}

func (s activateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 install step", slog.String("ID", s.Id()))
	return nil
}
