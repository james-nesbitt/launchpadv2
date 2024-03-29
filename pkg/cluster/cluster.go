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
}

// Validate the cluster configuration
func (c Cluster) Validate(ctx context.Context) error {

	// Dependency checking
	c.matchRequirements(ctx)

	if urs := c.requirements.UnMet(ctx); len(urs) > 0 {
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
func (c *Cluster) matchRequirements(ctx context.Context) error {
	c.requirements = dependency.Requirements{}

	// Collect all the things that provide dependencies
	fds := []dependency.FullfillsDependencies{}

	if cfd, ok := interface{}(c).(dependency.FullfillsDependencies); ok {
		fds = append(fds, cfd)
	}
	for _, h := range c.Hosts {
		if fd, ok := h.(dependency.FullfillsDependencies); ok {
			fds = append(fds, fd)
		}
	}
	for _, c := range c.Components {
		if fd, ok := c.(dependency.FullfillsDependencies); ok {
			fds = append(fds, fd)
		}
	}

	// Run through all source of dependencies, separately, so that we
	// Can do some logging and auditing
	for _, h := range c.Hosts {
		if hd, ok := h.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				if err := dependency.FullfillRequirements(ctx, r, fds); err != nil {
					slog.DebugContext(ctx, "host dependency requirement not met", slog.Any("cluster", c), slog.Any("requirement", r), slog.Any("host", h))
				}
				c.requirements = append(c.requirements, r)
			}
		}
	}
	for _, cc := range c.Components {
		if hd, ok := cc.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				if err := dependency.FullfillRequirements(ctx, r, fds); err != nil {
					slog.DebugContext(ctx, "component dependency requirement not met", slog.Any("cluster", c), slog.Any("requirement", r), slog.Any("component", cc))
				}
				c.requirements = append(c.requirements, r)
			}
		}
	}

	if ud := c.requirements.UnMet(ctx); len(ud) > 0 {
		return fmt.Errorf("%w; Unmet depedencies: %+v", ErrClusterDependenciesNotMet, ud)
	}
	return nil
}
