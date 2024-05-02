package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:msr4-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 discover step", slog.String("ID", s.Id()))
	return nil
}
