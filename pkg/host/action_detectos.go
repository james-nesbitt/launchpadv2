package host

import (
	"context"
	"fmt"
	"log/slog"
)

type detectOsStep struct {
	id string
}

func (s detectOsStep) Id() string {
	return fmt.Sprintf("%s:discover", s.id)
}

func (s detectOsStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running host detect-os step", slog.String("ID", s.id))
	return nil
}
