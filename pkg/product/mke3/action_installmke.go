package mke3

import (
	"context"
	"fmt"
	"log/slog"
)

type installMKEStep struct {
	id string
}

func (s installMKEStep) Id() string {
	return fmt.Sprintf("%s:mke3-install", s.id)
}

func (s installMKEStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE3 install step", slog.String("ID", s.id))
	return nil
}
