package dependency

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// MatchRequirements match requirements to dependencies, using sets of RequiresDependencies and ProvidesDependencies
//
// Each requirer's requirements will be passed to providers until a match is returned. The match is submitted to the requirement
// For completion, all rquirements and dpeendnecies are returned. You can check the requirers to see if any are unmatched.
func MatchRequirements(ctx context.Context, rds []RequiresDependencies, pds []ProvidesDependencies) (Requirements, Dependencies, error) {
	rs := Requirements{}
	ds := Dependencies{}

	slog.DebugContext(ctx, "Cluster: start component dependency matching")
	for _, rd := range rds {
		slog.DebugContext(ctx, "cluster: matching for dependency requirer", slog.Any("requirer", rd))
		for _, r := range rd.RequiresDependencies(ctx) {
			if d := r.Matched(ctx); d != nil {
				slog.DebugContext(ctx, "cluster: requirement already matched", slog.Any("cdependency", d), slog.Any("requirement", r))
				continue
			}

			rs = append(rs, r)

			if d, err := MatchRequirementDependency(ctx, r, pds); err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("cluster: requirement '%s' not matched: %s", r.Id(), err.Error()), slog.Any("requirement", r), slog.Any("error", err))
			} else if d == nil { // this should never happen, but nil err sometimes means no match
				slog.WarnContext(ctx, fmt.Sprintf("cluster: requirement '%s' not matched (empty)", r.Id()), slog.Any("requirement", r), slog.Any("error", err))
			} else if err := r.Match(d); err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("cluster: component requirement dependency match error %s->%s : %s", r.Id(), d.Id(), err.Error()), slog.Any("requirement", r), slog.Any("dependency", d), slog.Any("error", err))
			} else if d := r.Matched(ctx); d == nil {
				slog.WarnContext(ctx, fmt.Sprintf("cluster: component requirement dependency match failed (didn't stick) %s", r.Id()), slog.Any("requirement", r), slog.Any("dependency", d))
			} else {
				slog.DebugContext(ctx, "cluster: component requirement dependency matched", slog.Any("requirement", r), slog.Any("dependency", d))
				ds = append(ds, d)
			}
		}
	}
	slog.DebugContext(ctx, "Cluster: start component dependency matching")
	return rs, ds, nil
}

// MatchRequirementDependency Build a dependency for a requirement
//
//		Find a Handler which can handle the Requirement, and make it produce a Dependency.
//		Associate the Dependency with the Requirement, but return it as well so that it
//		can be collected.
//
//	    NOTE: WE DO NOT TELL THE REQUIREMENT ABOUT THE DEPENDENCY - YOU NEED TO DO THAT
func MatchRequirementDependency(ctx context.Context, r Requirement, pds []ProvidesDependencies) (Dependency, error) {
	if len(pds) == 0 {
		return nil, fmt.Errorf("%w; no dependency handlers providers for requirement %s", ErrNotHandled, r.Id())
	}

	for i, pd := range pds {
		if pd == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%d handler empty", i), slog.Any("requirement", r), slog.Any("handler", pd))
			continue // stupidity check
		}

		d, err := pd.ProvidesDependencies(ctx, r)

		// it is allowed to return no error, and no dependency,
		// which just means that the req wasn't handled
		if d == nil && err == nil {
			err = ErrNotHandled
		}

		if err == nil {
			slog.DebugContext(ctx, fmt.Sprintf("dependency: requirement '%s'match with %s", r.Id(), d.Id()), slog.Any("requirement", r), slog.Any("dependency", d), slog.Any("handler", pd))
			return d, nil // successful dependency creation
		}

		if errors.Is(err, ErrShouldHaveHandled) {
			slog.WarnContext(ctx, fmt.Sprintf("dependency: Dependency generation failure: (%s) %s", r.Id(), err.Error()), slog.Any("requirement", r), slog.Any("handler", pd), slog.String("error", err.Error()))
			return d, err // Handler says that it should have handled it, but couldn't, so return an error
		}

		if !errors.Is(err, ErrNotHandled) {
			slog.WarnContext(ctx, fmt.Sprintf("dependency: Unknown Dependency generation error: (%s) %s", r.Id(), err.Error()), slog.String("error", err.Error()))
		}
	}

	return nil, ErrNotHandled
}
