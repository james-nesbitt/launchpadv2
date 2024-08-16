package k0s

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

type uninstallK0sStep struct {
	baseStep
	id string
}

func (s uninstallK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-uninstall", s.id)
}

func (s uninstallK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s uninstall step", slog.String("ID", s.id))

	errs := []error{}

	if whs, whserr := s.c.GetAllHosts(ctx); whserr != nil {

	} else {
		if err := whs.Each(ctx, func(ctx context.Context, h *host.Host) error {
			slog.InfoContext(ctx, fmt.Sprintf("%s: k0s reset", h.Id()))
			return hostReset(ctx, h)
		}); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func hostReset(ctx context.Context, h *host.Host) error {
	kh := HostGetK0s(h)
	if kh == nil {
		return fmt.Errorf("%s: not a K0s host", h.Id())
	}

	kh.K0sStop(ctx)
	return kh.K0sReset(ctx)
}
