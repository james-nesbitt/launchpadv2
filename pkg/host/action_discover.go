package host

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	id string
}

func (ds discoverStep) Id() string {
	return fmt.Sprintf("%s:discover", ds.id)
}

func (ds discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running host discover step", slog.String("ID", ds.id))
	return nil
}
