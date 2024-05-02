package msr2

import (
	"context"
	"fmt"
	"log/slog"
)

type activateMSRStep struct {
	id string
}

func (s activateMSRStep) Id() string {
	return fmt.Sprintf("%s:msr2-activate", s.id)
}

func (s activateMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR2 activate step", slog.String("ID", s.Id()))
	return nil
}
