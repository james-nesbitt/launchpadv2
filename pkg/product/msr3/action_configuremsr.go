package msr3

import (
	"context"
	"fmt"
	"log/slog"
)

type configureMSRStep struct {
	id string
}

func (s configureMSRStep) Id() string {
	return fmt.Sprintf("%s:msr3-configure", s.id)
}

func (s configureMSRStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "RUNNING MSR3 CONFIGURE STEP", slog.String("ID", s.id))
	return nil
}
