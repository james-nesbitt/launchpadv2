package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type activateMSRStep struct {
	id string
}

func (s activateMSRStep) Id() string {
	return fmt.Sprintf("%s:msr4-activate", s.id)
}

func (s activateMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 activate step", slog.String("ID", s.Id()))
	return nil
}
