package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

type configureK0sStep struct {
	baseStep
	id string
}

func (s configureK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-configure", s.id)
}

func (s configureK0sStep) Validate(_ context.Context) error {
	return nil
}

func (s configureK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s configure step", slog.String("ID", s.Id()))

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return hserr
	}

	baseCfg := s.c.config.K0sConfig
	csans := s.c.CollectClusterSans(ctx)

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		kh := HostGetK0s(h)

		if err := kh.WriteK0sConfig(ctx, baseCfg, csans); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
