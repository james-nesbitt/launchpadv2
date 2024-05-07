package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type swarmActivateStep struct {
	baseStep

	id string
}

func (s swarmActivateStep) Id() string {
	return fmt.Sprintf("%s:mcr-swarm-activate", s.id)
}

func (s swarmActivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR swarm-activate step", slog.String("ID", s.Id()))
	return nil
}
