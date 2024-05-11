package cluster

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

var (
	ErrClusterValidationFailed   = errors.New("cluster validation failed")
	ErrClusterDependenciesNotMet = errors.New("cluster dependencies not met")
)

// Cluster function handler for the complete cluster.
type Cluster struct {
	Components component.Components
}

// Validate the cluster configuration.
func (cl *Cluster) Validate(ctx context.Context) error {
	cerrs := []error{}

	// Dependency checking
	_, _, merr := cl.matchRequirements(ctx)
	if merr != nil {
		return fmt.Errorf("cluster validate failed on matching requirements: %w", merr)
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

// Command build a command using the dependencies and components from the cluster.
func (cl *Cluster) Command(ctx context.Context, key string) (*action.Command, error) {
	slog.InfoContext(ctx, fmt.Sprintf("%s Command build", key))

	cmd := action.NewEmptyCommand(key)

	// Dependency checking
	_, ds, merr := cl.matchRequirements(ctx)
	if merr != nil {
		return cmd, fmt.Errorf("cluster validate failed on matching requirements: %w", merr)
	}

	// add all cluster dependencies to the command
	cmd.Dependencies.Merge(ds)

	errs := []error{}
	for _, c := range cl.Components {
		cmdb, ok := c.(action.CommandHandler)
		if !ok {
			slog.WarnContext(ctx, "component is not a command builder", slog.Any("component", c))
			continue
		}

		if err := cmdb.CommandBuild(ctx, cmd); err != nil {
			slog.WarnContext(ctx, "component command build error", slog.Any("component", c), slog.Any("error", err))
			errs = append(errs, err)
		} else {
			slog.InfoContext(ctx, "component command build", slog.Any("component", c))
		}
	}

	if len(errs) > 0 {
		return cmd, fmt.Errorf("failed to build command: %s", errors.Join(errs...).Error())
	}
	return cmd, nil
}

// match all requirements from cluster components with dependency fullfillers.
//
//	This allows the list of requirements to be analyzed to make sure that
//	all of the cluster dependencies are met.
func (cl *Cluster) matchRequirements(ctx context.Context) (dependency.Requirements, dependency.Dependencies, error) {
	// Collect all the things that provide dependencies
	rds := []dependency.RequiresDependencies{}
	pds := []dependency.ProvidesDependencies{}
	for _, c := range cl.Components {
		if rd, ok := c.(dependency.RequiresDependencies); ok {
			slog.DebugContext(ctx, fmt.Sprintf("including component '%s' as a dependency requirer", c.Name()), slog.Any("component", c))
			rds = append(rds, rd)
		}
		if pd, ok := c.(dependency.ProvidesDependencies); ok {
			slog.DebugContext(ctx, fmt.Sprintf("including component '%s' as a dependency provider", c.Name()), slog.Any("component", c))
			pds = append(pds, pd)
		}
	}

	slog.DebugContext(ctx, "Cluster: start component dependency matching")
	rs, ds, merr := dependency.MatchRequirements(ctx, rds, pds)
	slog.DebugContext(ctx, "Cluster: finished component dependency matching")
	if merr != nil {
		return rs, ds, fmt.Errorf("cluster dependency matching failed: %s", merr.Error())
	}

	if urs := rs.UnMatched(ctx); len(urs) > 0 {
		for _, ur := range urs {
			slog.ErrorContext(ctx, "UnMatched cluster dependency requirement", slog.Any("requirement", ur))
		}
		return rs, ds, fmt.Errorf("%w; Unmet dependencies: %+v", ErrClusterDependenciesNotMet, urs)
	}
	return rs, ds, nil
}
