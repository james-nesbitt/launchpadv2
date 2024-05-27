package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

/**
 * K0S Installation
 *
 * Here we get the K0S binaries onto the machines, install
 * Any services need, and install any config needed.
 *
 * Running the services and connecting them is in the activation
 * phase.
 */

type installK0sStep struct {
	baseStep
	id    string
	Force bool
}

func (s installK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-install", s.id)
}

func (s installK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s install step", slog.String("ID", s.id))

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return hserr
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.DebugContext(ctx, fmt.Sprintf("%s: downloading k0s binary on host", h.Id()))

		kh := HostGetK0s(h)

		if s.c.config.ShouldDownload() {
			if err := kh.DownloadK0sBinary(ctx, s.c.config.Version); err != nil {
				return err
			}
		} else {
			if err := kh.UploadK0sBinary(ctx, s.c.config.Version); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
