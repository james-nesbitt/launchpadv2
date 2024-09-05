package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type downloadClientBundleStep struct {
	baseStep
	id string
}

func (s downloadClientBundleStep) Id() string {
	return fmt.Sprintf("%s:mke4-client-bundle", s.id)
}

func (s downloadClientBundleStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 client-bundle generate step", slog.String("ID", s.Id()))
	return nil
}
