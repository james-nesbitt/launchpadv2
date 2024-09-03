package k0s

import (
	"context"
	"fmt"
	"log/slog"
)

type kubernetesConfStep struct {
	baseStep
	id string
}

func (s kubernetesConfStep) Id() string {
	return fmt.Sprintf("%s:k0s-kube-config", s.id)
}

func (s kubernetesConfStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s kubeconfig step", slog.String("ID", s.id))
	return nil
}
