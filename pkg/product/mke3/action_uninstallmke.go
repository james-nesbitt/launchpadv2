package mke3

import (
	"context"
	"fmt"
	"log/slog"
)

type uninstallMKEStep struct {
	id string
}

func (s uninstallMKEStep) Id() string {
	return fmt.Sprintf("%s:mke3-uninstall", s.id)
}

func (s uninstallMKEStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE3 uninstall step", slog.String("ID", s.id))
	return nil
}
