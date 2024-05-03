package mcr

import (
	"context"
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

	if err := hs.Each(ctx, discoverVersion); err != nil {
		slog.WarnContext(ctx, "MCR discovery failed to discover installation")
		if s.failOnError {
			return err
		}
	}
	return nil
}

func discoverVersion(ctx context.Context, h *dockerhost.Host) error {
	v, err := h.Docker(ctx).ServerVersion(ctx)
	if err != nil {
		return fmt.Errorf("%s: version discover error: %s", h.Id(), err)
	}

	slog.InfoContext(ctx, "host version discovered", slog.Any("version", v))

	return nil
}
