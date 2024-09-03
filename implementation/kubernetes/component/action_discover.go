package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:k8s-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running kubernetes discover step", slog.String("ID", s.id))

	ki, kierr := s.c.kubernetesImplementation(ctx)
	if kierr != nil {
		return kierr
	}

	ksv, ksverr := ki.ServerVersion(ctx)
	if ksverr != nil {
		return ksverr
	}

	slog.DebugContext(ctx, "kubernetes server discovery successful", slog.Any("ServerVersion", ksv))
	return nil
}
