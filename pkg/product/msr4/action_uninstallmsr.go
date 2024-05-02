package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallMSRStep struct {
	id string
}

func (s uninstallMSRStep) Id() string {
	return fmt.Sprintf("%s:msr4-uninstall", s.id)
}

func (s uninstallMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 uninstall step", slog.String("ID", s.Id()))
	return nil
}
