package mcr

import (
	"context"
	"fmt"
	"log/slog"
)

type prepareMCRNodesStep struct {
	id string
}

func (s prepareMCRNodesStep) Id() string {
	return fmt.Sprintf("%s:mcr-configure", s.id)
}

func (s prepareMCRNodesStep) Validate(_ context.Context) error {
	return nil
}

func (s prepareMCRNodesStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR prepare-nodes step", slog.String("ID", s.id))
	return nil
}
