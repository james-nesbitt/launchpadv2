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
	baseStep

	id string
	d  dependency.Dependency
}

func (ds discoverStep) Id() string {
	return fmt.Sprintf("%s:%s", DiscoverStepId, ds.id)
}

func (ds discoverStep) Validate(ctx context.Context) error {
	hd, ok := ds.d.(HostsDependency)
	if !ok {
		return fmt.Errorf("%w; %s has a Dependency of the wrong type", dependency.ErrDependencyNotMatched, ds.Id())
	}

	hs, err := hd.ProduceHosts(ctx)
	if err != nil {
		return fmt.Errorf("%w; %s failed to1 produce hosts: %s", dependency.ErrDependencyNotMatched, ds.Id(), err.Error())
	}
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

	hs, err := hd.ProduceHosts(ctx)
	if err != nil {
		return fmt.Errorf("%w; %s failed to produce hosts: %s", dependency.ErrDependencyNotMatched, ds.Id(), err.Error())
	}
	if len(hs) == 0 {
		slog.WarnContext(ctx, "HostComponent:DiscoverStep dependency has no hosts.")
		return nil
	}

	cerr := hs.Each(ctx, func(ctx context.Context, h *Host) error {
		hd := HostGetDiscover(h)
		if hd == nil {
			return fmt.Errorf("%s host discover impossible. No host plugin registered that performs discovery", h.Id())
		}

		if err := hd.Discover(ctx); err != nil {
			return fmt.Errorf("%s host discover failed: %s", h.Id(), err.Error())
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: host discover successful", h.Id()), slog.Any("host", h))
		return nil
	})

	return cerr
}
