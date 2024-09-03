package msr2

import (
	"context"
	"fmt"
	"log/slog"
)

type configureMSRStep struct {
	id string
}

func (s configureMSRStep) Id() string {
	return fmt.Sprintf("%s:msr2-configure", s.id)
}

func (s configureMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR2 configure step", slog.String("ID", s.Id()))
	return nil
}
