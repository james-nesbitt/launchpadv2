package msr3

import (
	"context"
	"fmt"
	"log/slog"
)

type installMSRStep struct {
	id string
}

func (s installMSRStep) Id() string {
	return fmt.Sprintf("%s:msr3-install", s.id)
}

func (s installMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "RUNNING MSR INSTALL STEP", slog.String("ID", s.Id()))
	return nil
}
