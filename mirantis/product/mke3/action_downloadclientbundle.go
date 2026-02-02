package mke3

import (
	"context"
	"fmt"
	"log/slog"
)

type downloadClientBundleStep struct {
	baseStep
	id string
}

func (s downloadClientBundleStep) ID() string {
	return fmt.Sprintf("%s:mke3-client-bundle", s.id)
}

func (s downloadClientBundleStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE3 client-bundle generate step", slog.String("ID", s.ID()))
	return nil
}
