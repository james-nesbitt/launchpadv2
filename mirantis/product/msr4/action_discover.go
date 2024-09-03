package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:msr4-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 discover step", slog.String("ID", s.Id()))

	r, rerr := s.c.getRelease(ctx)
	if rerr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("MSR4 release %s not found", s.c.Name()))
		return rerr
	}

	slog.InfoContext(ctx, fmt.Sprintf("MSR4 release %s found", r.Name), slog.Any("release.info", r.Info))
	return nil
}
