package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type configureMSRStep struct {
	id string
}

func (s configureMSRStep) Id() string {
	return fmt.Sprintf("%s:msr4-configure", s.id)
}

func (s configureMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 configure step", slog.String("ID", s.Id()))
	return nil
}
