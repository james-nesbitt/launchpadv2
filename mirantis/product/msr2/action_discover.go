package msr2

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:msr2-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR2 discover step", slog.String("ID", s.Id()))
	return nil
}
