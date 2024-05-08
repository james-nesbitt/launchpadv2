package mcr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

type discoverStep struct {
	baseStep

	id          string
	failOnError bool
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:mcr-discover", s.id)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR discover step", slog.String("ID", s.Id()))

	hs, gherr := s.c.GetAllHosts(ctx)
	if gherr != nil {
		return fmt.Errorf("could not discover hosts: %s", gherr.Error())
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		hs := s.c.state.hosts.GetOrCreate(h.Id())

		verr := s.discoverVersion(ctx, h)
		ierr := s.discoverInfo(ctx, h)

		if ierr != nil || verr != nil {
			return errors.Join(ierr, verr)
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: state discovered", h.Id()), slog.Any("host", h), slog.Any("version", hs.version), slog.Any("info", hs.info))

		return nil
	}); err != nil {
		slog.WarnContext(ctx, "MCR discovery failed to discover installation", slog.Any("error", err.Error()))
		if s.failOnError {
			return err
		}
	}
	return nil
}

func (s discoverStep) discoverVersion(ctx context.Context, h *dockerhost.Host) error {
	v, err := h.Docker(ctx).Version(ctx)
	if err != nil {
		return fmt.Errorf("%s: version discover error: %s", h.Id(), err)
	}

	sh := s.c.state.hosts.GetOrCreate(h.Id())
	sh.setVersion(v)

	return nil
}

func (s discoverStep) discoverInfo(ctx context.Context, h *dockerhost.Host) error {
	i, err := h.Docker(ctx).Info(ctx)
	if err != nil {
		return fmt.Errorf("%s: info discover error: %s", h.Id(), err)
	}

	sh := s.c.state.hosts.GetOrCreate(h.Id())
	sh.setInfo(i)

	return nil
}
