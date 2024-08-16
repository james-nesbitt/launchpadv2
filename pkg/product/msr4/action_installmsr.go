package msr4

import (
	"context"
	"fmt"
	"log/slog"
)

type installMSRStep struct {
	baseStep
	id string
}

func (s installMSRStep) Id() string {
	return fmt.Sprintf("%s:msr4-install", s.id)
}

func (s installMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MSR4 install step", slog.String("ID", s.Id()))

	if err := s.c.addMsr4HelmRepo(ctx); err != nil {
		return err
	}

	return nil
}
