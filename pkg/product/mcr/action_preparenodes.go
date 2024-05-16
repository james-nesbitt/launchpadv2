package mcr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

var (
	basePackages = []string{"curl"}
)

type prepareMCRNodesStep struct {
	baseStep

	id string
}

func (s prepareMCRNodesStep) Id() string {
	return fmt.Sprintf("%s:mcr-prepare-nodes", s.id)
}

func (s prepareMCRNodesStep) Run(ctx context.Context) error {
	hs, hsgerr := s.c.GetAllHosts(ctx)
	if hsgerr != nil {
		return hsgerr
	}

	if err := hs.Each(ctx, installBasePackages); err != nil {
		return err
	}
	return nil
}

func installBasePackages(ctx context.Context, h *host.Host) error {
	he := exec.HostGetExecutor(h)
	hp := exec.HostGetPlatform(h)

	if hp.IsWindows(ctx) {
		return nil
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: installing base packages", h.Id()), slog.Any("packages", basePackages))
	return he.InstallPackages(ctx, basePackages)
}
