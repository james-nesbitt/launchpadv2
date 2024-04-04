package cluster

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

var (
	ErrClusterDependenciesNotMet = errors.New("cluster dependencies not met")
)

// Cluster function handler for the complete cluster
type Cluster struct {
	Hosts      host.Hosts
	Components component.Components

	requirements dependency.Requirements
	dependencies dependency.Dependencies
}

// Validate the cluster configuration
func (cl Cluster) Validate(ctx context.Context) error {

	// Dependency checking
	if err := cl.matchRequirements(ctx); err != nil {
		return fmt.Errorf("cluster validate failed on matching requirements: %w", err)
	}

	if urs := cl.requirements.UnMatched(ctx); len(urs) > 0 {
		var urses []string
		for _, ur := range urs {
			urses = append(urses, ur.Describe())
		}

		return fmt.Errorf("%w; %s", ErrClusterDependenciesNotMet, strings.Join(urses, "\n"))
	}

	return nil
}

// match all requirements from cluster components with dependency fullfillers.
//
//	This allows the list of requirements to be analyzed to make sure that
//	all of the cluster dependencies are met.
func (cl *Cluster) matchRequirements(ctx context.Context) error {
	cl.requirements = dependency.Requirements{}
	cl.dependencies = dependency.Dependencies{}

	// Collect all the things that provide dependencies
	fds := []dependency.FullfillsDependencies{}

	if cfd, ok := interface{}(cl).(dependency.FullfillsDependencies); ok {
		fds = append(fds, cfd)
	}
	for _, h := range cl.Hosts {
		if fd, ok := h.(dependency.FullfillsDependencies); ok {
			fds = append(fds, fd)
		}
	}
	for _, c := range cl.Components {
		if fd, ok := c.(dependency.FullfillsDependencies); ok {
			fds = append(fds, fd)
		}
	}

	// Run through all source of dependencies, separately, so that we
	// Can do some logging and auditing
	for _, h := range cl.Hosts {
		if hd, ok := h.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				if d, err := dependency.MatchRequirements(ctx, r, fds); err != nil {
					slog.DebugContext(ctx, "host dependency requirement not met", slog.Any("cluster", cl), slog.Any("requirement", r), slog.Any("host", h))
				} else {
					r.Match(d)
					slog.DebugContext(ctx, "host dependency requirement met", slog.Any("cluster", cl), slog.Any("requirement", r), slog.Any("host", h))
					cl.requirements = append(cl.requirements, r)
					cl.dependencies = append(cl.dependencies, d)
				}
			}
		}
	}
	for _, cc := range cl.Components {
		if hd, ok := cc.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				slog.Error("Checking Component for provides: %+v", slog.Any("requirement", r))
				if d, err := dependency.MatchRequirements(ctx, r, fds); err != nil {
					slog.DebugContext(ctx, "component dependency requirement not met", slog.Any("cluster", cl), slog.Any("requirement", r), slog.Any("component", cc))
				} else {
					r.Match(d)
					slog.DebugContext(ctx, "component dependency requirement met and matched", slog.Any("cluster", cl), slog.Any("requirement", r), slog.Any("dependency", d), slog.Any("component", cc))
					cl.requirements = append(cl.requirements, r)
					cl.dependencies = append(cl.dependencies, d)
				}
			}
		}
	}

	if ud := cl.requirements.UnMatched(ctx); len(ud) > 0 {
		return fmt.Errorf("%w; Unmet depedencies: %+v :: \n %+v", ErrClusterDependenciesNotMet, ud, fds)
	}
	return nil
}
