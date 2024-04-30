package host

import (
	"context"
	"fmt"
	"log/slog"
)

type runHooksStep struct {
	id string
}

func (s runHooksStep) Id() string {
	return fmt.Sprintf("%s:discover", s.id)
}

func (s runHooksStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running host hooks", slog.String("ID", s.id))
	return nil
}
