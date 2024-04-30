package host

import (
	"context"
	"fmt"
	"log/slog"
)

type disconnectStep struct {
	id string
}

func (s disconnectStep) Id() string {
	return fmt.Sprintf("%s:discover", s.id)
}

func (s disconnectStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running host disconnect step", slog.String("ID", s.id))
	return nil
}
