package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) ID() string {
	return fmt.Sprintf("%s:k0s-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s discover step", slog.String("ID", s.id))

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return hserr
	}

	return hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		kh := HostGetK0s(h)

		if v, err := kh.Version(ctx); err == nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: discovered k0s version %s", h.ID(), v.K0s))
		} else {
			slog.DebugContext(ctx, fmt.Sprintf("%s: k0s version not found", h.ID()))
		}

		if st, err := kh.Status(ctx); err == nil {
			slog.InfoContext(ctx, fmt.Sprintf("%s: discovered k0s status (role=%s, workloads=%t)", h.ID(), st.Role, st.Workloads))
		} else {
			slog.DebugContext(ctx, fmt.Sprintf("%s: k0s status not found", h.ID()))
		}

		return nil
	})
}
