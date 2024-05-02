package mke4

import (
	"context"
	"fmt"
	"log/slog"
)

type prepareNodesStep struct {
	id string
}

func (s prepareNodesStep) Id() string {
	return fmt.Sprintf("%s:mke4-prepare-nodes", s.id)
}

func (s prepareNodesStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MKE4 prepare-nodes step", slog.String("ID", s.Id()))
	return nil
}
