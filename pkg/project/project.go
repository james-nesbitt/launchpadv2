package project

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
	ErrProjectValidationFailed   = errors.New("project validation failed")
	ErrProjectDependenciesNotMet = errors.New("project dependencies not met")
)

// Project function handler for the complete project.
type Project struct {
	Components component.Components
}

// Validate the project configuration.
func (p *Project) Validate(ctx context.Context) error {
	cerrs := []error{}

	// Dependency checking
	_, _, merr := p.matchRequirements(ctx)
	if merr != nil {
		return fmt.Errorf("project validate failed on matching requirements: %w", merr)
	}

	// Validate components
	verrs := []error{}
	for _, c := range p.Components {
		if err := c.Validate(ctx); err != nil {
			slog.WarnContext(ctx, "Component validation failure", slog.Any("component", c))
			verrs = append(verrs, err)
		}
	}
	if len(verrs) > 0 {
		cerrs = append(cerrs, fmt.Errorf("%w: %s", ErrProjectValidationFailed, errors.Join(verrs...).Error()))
	}

	if len(cerrs) > 0 {
		return fmt.Errorf("%w; component validation errors: \n %s", ErrProjectValidationFailed, errors.Join(cerrs...).Error())
	}

	return nil
}

// Command build a command using the dependencies and components from the project.
func (p *Project) Command(ctx context.Context, key string) (*action.Command, error) {
	slog.InfoContext(ctx, fmt.Sprintf("%s Command build", key))

	cmd := action.NewEmptyCommand(key)

	// Dependency checking
	_, ds, merr := p.matchRequirements(ctx)
	if merr != nil {
		return cmd, fmt.Errorf("project validate failed on matching requirements: %w", merr)
	}

	// add all project dependencies to the command
	cmd.Dependencies.Merge(ds)

	errs := []error{}
	for _, c := range p.Components {
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

// match all requirements from project components with dependency fullfillers.
//
//	This allows the list of requirements to be analyzed to make sure that
//	all of the project dependencies are met.
func (p *Project) matchRequirements(ctx context.Context) (dependency.Requirements, dependency.Dependencies, error) {
	// Collect all the things that provide dependencies
	rds := []dependency.RequiresDependencies{}
	pds := []dependency.ProvidesDependencies{}
	for _, c := range p.Components {
		if rd, ok := c.(dependency.RequiresDependencies); ok {
			slog.DebugContext(ctx, fmt.Sprintf("including component '%s' as a dependency requirer", c.Name()), slog.Any("component", c))
			rds = append(rds, rd)
		}
		if pd, ok := c.(dependency.ProvidesDependencies); ok {
			slog.DebugContext(ctx, fmt.Sprintf("including component '%s' as a dependency provider", c.Name()), slog.Any("component", c))
			pds = append(pds, pd)
		}
	}

	slog.DebugContext(ctx, "Project: start component dependency matching")
	rs, ds, merr := dependency.MatchRequirements(ctx, rds, pds)
	slog.DebugContext(ctx, "Project: finished component dependency matching")
	if merr != nil {
		return rs, ds, fmt.Errorf("project dependency matching failed: %s", merr.Error())
	}

	if urs := rs.UnMatched(ctx); len(urs) > 0 {
		for _, ur := range urs {
			slog.ErrorContext(ctx, "UnMatched project dependency requirement", slog.Any("requirement", ur))
		}
		return rs, ds, fmt.Errorf("%w; Unmet dependencies: %+v", ErrProjectDependenciesNotMet, urs)
	}
	return rs, ds, nil
}
