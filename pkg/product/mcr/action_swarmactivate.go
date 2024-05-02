package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type swarmActivateStep struct {
	id string
}

func (s swarmActivateStep) Id() string {
	return fmt.Sprintf("%s:mcr-swarm-activate", s.id)
}

func (s swarmActivateStep) Validate(_ context.Context) error {
	return nil
}

func (s swarmActivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR swarm-activate step", slog.String("ID", s.Id()))
	return nil
}
