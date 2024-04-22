package msr2

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallMSRStep struct {
	id string
}

func (s uninstallMSRStep) Id() string {
	return fmt.Sprintf("%s:msr2-uninstall", s.id)
}

func (s uninstallMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR2 uninstall step", slog.String("ID", s.id))
	return nil
}
