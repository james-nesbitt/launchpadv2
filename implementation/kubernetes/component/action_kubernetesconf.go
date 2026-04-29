package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type kubernetesConfStep struct {
	baseStep
	id string
}

func (s kubernetesConfStep) ID() string {
	return fmt.Sprintf("%s:kubernetes-conf", s.id)
}

// Run exports the cluster kubeconfig to stdout.
//
// This step is run as part of the KubeConf command to retrieve the kubeconfig
// for the cluster managed by this component.
func (s kubernetesConfStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running kubeconfig step", slog.String("ID", s.id))

	cb := s.c.config.KubeConfig()
	if len(cb) == 0 {
		return fmt.Errorf("%s: kubeconfig is empty; ensure the cluster is configured and reachable", s.id)
	}

	// Write the kubeconfig to stdout for consumption by callers
	// (e.g., piped to kubectl or saved to a file by the shell).
	if _, err := os.Stdout.Write(cb); err != nil {
		return fmt.Errorf("%s: failed to write kubeconfig to stdout: %w", s.id, err)
	}

	return nil
}
