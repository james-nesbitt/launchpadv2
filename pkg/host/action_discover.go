package host

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	DiscoverStepId = "HostComponent:DiscoverStep"
)

type discoverStep struct {
	id string

	d dependency.Dependency
}

func (ds discoverStep) Id() string {
	return fmt.Sprintf("%s:%s", DiscoverStepId, ds.id)
}

func (ds discoverStep) Validate(ctx context.Context) error {
	hd, ok := ds.d.(HostsDependency)
	if !ok {
		return fmt.Errorf("%w; %s has a Dependency of the wrong type", dependency.ErrDependencyNotMatched, ds.Id())
	}

	hs := hd.ProduceHosts(ctx)
	if len(hs) == 0 {
		slog.WarnContext(ctx, "HostComponent:DiscoverStep dependency has no hosts.")
		return nil
	}

	return nil
}

func (ds discoverStep) Run(ctx context.Context) error {
	hd, ok := ds.d.(HostsDependency)
	if !ok {
		return fmt.Errorf("%w; %s has a Dependency of the wrong type", dependency.ErrDependencyNotMatched, ds.Id())
	}

	hs := hd.ProduceHosts(ctx)
	if len(hs) == 0 {
		slog.WarnContext(ctx, "HostComponent:DiscoverStep dependency has no hosts.")
		return nil
	}

	cerr := hs.Each(ctx, func(ctx context.Context, h Host) error {
		if err := h.Connect(ctx); err != nil {
			return fmt.Errorf("%s host connection failed: %s", h.Id(), err.Error())
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: host connection successful", h.Id()), slog.Any("host", h))
		return nil
	})

	return cerr
}
