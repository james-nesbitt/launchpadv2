package msr2

import (
	"context"
	"fmt"
	"log/slog"
)

type installMSRStep struct {
	id string
}

func (s installMSRStep) Id() string {
	return fmt.Sprintf("%s:msr2-install", s.id)
}

func (s installMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR2 install step", slog.String("ID", s.id))
	return nil
}
