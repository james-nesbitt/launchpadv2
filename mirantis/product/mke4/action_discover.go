package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:mke4-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 discover step", slog.String("ID", s.Id()))
	return nil
}
