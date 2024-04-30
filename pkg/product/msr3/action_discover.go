package msr3

import (
	"context"
	"fmt"
	"log/slog"
)

// Combine GatherFacts, DetectOS and

type discoverStep struct {
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:msr3-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "RUNNING MSR3 DISCOVER STEP", slog.String("ID", s.id))
	return nil
}
