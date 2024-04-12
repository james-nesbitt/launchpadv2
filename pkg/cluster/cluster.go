package cluster

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

var (
	ErrClusterValidationFailed   = errors.New("cluster validation failed")
	ErrClusterDependenciesNotMet = errors.New("cluster dependencies not met")
)

// Cluster function handler for the complete cluster
type Cluster struct {
	Components component.Components

	requirements dependency.Requirements
	dependencies dependency.Dependencies
}

// Validate the cluster configuration
func (cl *Cluster) Validate(ctx context.Context) error {
	cerrs := []error{}

	// Dependency checking
	if err := cl.matchRequirements(ctx); err != nil {
		return fmt.Errorf("cluster validate failed on matching requirements: %w", err)
	} else if urs := cl.requirements.UnMatched(ctx); len(urs) > 0 {
		var urses []string
		for _, ur := range urs {
			urses = append(urses, ur.Describe())
		}

		slog.WarnContext(ctx, "Component dependencies unmatched", slog.Any("unmatched-dependencies", urs))
		cerrs = append(cerrs, fmt.Errorf("%w; %s", ErrClusterDependenciesNotMet, strings.Join(urses, "\n")))
	}

	// Validate components
	verrs := []error{}
	for _, c := range cl.Components {
		if err := c.Validate(ctx); err != nil {
			slog.WarnContext(ctx, "Component validation failure", slog.Any("component", c))
			verrs = append(verrs, err)
		}
	}
	if len(verrs) > 0 {
		cerrs = append(cerrs, fmt.Errorf("%w: %s", ErrClusterValidationFailed, errors.Join(verrs...).Error()))
	}

	if len(cerrs) > 0 {
		return fmt.Errorf("%w; component validation errors: \n %s", ErrClusterValidationFailed, errors.Join(cerrs...).Error())
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

	for _, c := range cl.Components {
		if fd, ok := c.(dependency.FullfillsDependencies); ok {
			slog.DebugContext(ctx, fmt.Sprintf("including component '%s' as a dependency provider", c.Name()), slog.Any("component", c))
			fds = append(fds, fd)
		}
	}

	for _, cc := range cl.Components {
		if hd, ok := cc.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				cl.requirements = append(cl.requirements, r)

				if d, err := dependency.GetRequirementDependency(ctx, r, fds); err != nil {
					slog.DebugContext(ctx, fmt.Sprintf("cluster: requirement '%s' not matched: %s", r.Id(), err.Error()), slog.Any("requirement", r), slog.Any("component", cc), slog.Any("error", err))
				} else if d == nil { // this should never happen, but nil err sometimes means no match
					slog.DebugContext(ctx, fmt.Sprintf("cluster: requirement '%s' not matched (empty): %s", r.Id(), err.Error()), slog.Any("requirement", r), slog.Any("component", cc), slog.Any("error", err))
				} else if err := r.Match(d); err != nil {
					slog.DebugContext(ctx, fmt.Sprintf("cluster: component requirement dependency match error %s->%s : %s", r.Id(), d.Id(), err.Error()), slog.Any("requirement", r), slog.Any("dependency", d), slog.Any("error", err))
				} else if d := r.Matched(ctx); d == nil {
					slog.DebugContext(ctx, fmt.Sprintf("cluster: component requirement dependency match failed (didn't stick) %s->%s", r.Id(), d.Id()), slog.Any("requirement", r), slog.Any("dependency", d), slog.Any("error", err))
				} else {
					slog.DebugContext(ctx, "cluster: component requirement dependency matched", slog.Any("requirement", r), slog.Any("dependency", d))
					cl.dependencies = append(cl.dependencies, d)
				}
			}
		}
	}

	if urs := cl.requirements.UnMatched(ctx); len(urs) > 0 {
		for _, ur := range urs {
			slog.ErrorContext(ctx, "UnMatched cluster dependency requirement", slog.Any("requirement", ur))
		}
		return fmt.Errorf("%w; Unmet depedencies: %+vv", ErrClusterDependenciesNotMet, urs)
	}
	return nil
}
