package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallOperatorStep struct {
	baseStep
	id string
}

func (s uninstallOperatorStep) Id() string {
	return fmt.Sprintf("%s:mke4-uninstall-operator", s.id)
}

func (s uninstallOperatorStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 operator uninstall step", slog.String("ID", s.Id()))
	return nil
}
