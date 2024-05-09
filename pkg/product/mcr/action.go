package mcr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

type baseStep struct {
	c MCR
}

func (s baseStep) Validate(ctx context.Context) error {
	if _, err := s.c.GetAllHosts(ctx); err != nil {
		return fmt.Errorf("MCR Step validation fail: Host validation: %s", err.Error())
	}

	return nil
}

// updateDockerState update the state with the latest docker information
func (s baseStep) updateDockerState(ctx context.Context) error {

	hs, gherr := s.c.GetAllHosts(ctx)
	if gherr != nil {
		return fmt.Errorf("could not discover hosts: %s", gherr.Error())
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: discovering MCR state", h.Id()), slog.Any("host", h))

		hs := s.c.state.hosts.GetOrCreate(h.Id())

		slog.DebugContext(ctx, fmt.Sprintf("%s: discovering MCR version", h.Id()), slog.Any("host", h))
		verr := s.discoverVersion(ctx, h)
		slog.DebugContext(ctx, fmt.Sprintf("%s: discovering MCR info", h.Id()), slog.Any("host", h))
		ierr := s.discoverInfo(ctx, h)

		if ierr != nil || verr != nil {
			return errors.Join(ierr, verr)
		}

		slog.DebugContext(ctx, fmt.Sprintf("%s: MCR state discovered", h.Id()), slog.Any("host", h), slog.Any("version", hs.version), slog.Any("info", hs.info))

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s baseStep) discoverVersion(ctx context.Context, h *dockerhost.Host) error {
	v, err := h.Docker(ctx).Version(ctx)
	if err != nil {
		return fmt.Errorf("%s: MCR version discover error: %s", h.Id(), err)
	}

	sh := s.c.state.hosts.GetOrCreate(h.Id())
	sh.setVersion(v)

	return nil
}

func (s baseStep) discoverInfo(ctx context.Context, h *dockerhost.Host) error {
	i, err := h.Docker(ctx).Info(ctx)
	if err != nil {
		return fmt.Errorf("%s: MCR info discover error: %s", h.Id(), err)
	}

	sh := s.c.state.hosts.GetOrCreate(h.Id())
	sh.setInfo(i)

	return nil
}
