package msr3

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallMSRStep struct {
	id string
}

func (s uninstallMSRStep) Id() string {
	return fmt.Sprintf("%s:msr-uninstall", s.id)
}

func (s uninstallMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "RUNNING MSR3 UNINSTALL STEP", slog.String("ID", s.id))
	return nil
}
