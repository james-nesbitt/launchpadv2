package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type activateMSRStep struct {
	baseStep
	id string
}

func (s activateMSRStep) Id() string {
	return fmt.Sprintf("%s:msr4-activate", s.id)
}

func (s activateMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 activate step", slog.String("ID", s.Id()))

	r, rerr := s.c.installMsr4Release(ctx)
	if rerr != nil {
		return rerr
	}

	slog.InfoContext(ctx, fmt.Sprintf("MSR4 release %s found", r.Name), slog.Any("release.info", r.Info))

	return nil
}
