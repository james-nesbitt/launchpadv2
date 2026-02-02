package host

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	DiscoverStepID = "HostComponent:DiscoverStep"
)

type discoverStep struct {
	baseStep

	id string
	d  dependency.Dependency
}

func (ds discoverStep) ID() string {
	return fmt.Sprintf("%s:%s", DiscoverStepID, ds.id)
}

func (ds discoverStep) Validate(ctx context.Context) error {
	hd, ok := ds.d.(HostsDependency)
	if !ok {
		return fmt.Errorf("%w; %s has a Dependency of the wrong type", dependency.ErrDependencyNotMatched, ds.ID())
	}

	hs, err := hd.ProduceHosts(ctx)
	if err != nil {
		return fmt.Errorf("%w; %s failed to1 produce hosts: %s", dependency.ErrDependencyNotMatched, ds.ID(), err.Error())
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
		return fmt.Errorf("%w; %s has a Dependency of the wrong type", dependency.ErrDependencyNotMatched, ds.ID())
	}

	hs, err := hd.ProduceHosts(ctx)
	if err != nil {
		return fmt.Errorf("%w; %s failed to produce hosts: %s", dependency.ErrDependencyNotMatched, ds.ID(), err.Error())
	}
	if len(hs) == 0 {
		slog.WarnContext(ctx, "HostComponent:DiscoverStep dependency has no hosts.")
		return nil
	}

	cerr := hs.Each(ctx, func(ctx context.Context, h *Host) error {
		hd := HostGetDiscover(h)
		if hd == nil {
			return fmt.Errorf("%s host discover impossible. No host plugin registered that performs discovery", h.ID())
		}

		slog.DebugContext(ctx, fmt.Sprintf("%s: Discovering host", h.ID()))
		if err := hd.Discover(ctx); err != nil {
			return fmt.Errorf("%s host discover failed: %s", h.ID(), err.Error())
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: host discover successful", h.ID()), slog.Any("host", h))
		return nil
	})

	return cerr
}
