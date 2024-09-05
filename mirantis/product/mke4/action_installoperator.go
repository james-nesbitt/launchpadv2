package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type installOperatorStep struct {
	baseStep
	id string
}

func (s installOperatorStep) Id() string {
	return fmt.Sprintf("%s:install-operator", s.id)
}

func (s installOperatorStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 Operator install step", slog.String("ID", s.Id()))
	return nil
}
