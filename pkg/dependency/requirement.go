package dependency

import (
	"context"
	"log/slog"
)

// HasDependencies can determine its dependencies
type HasDependencies interface {
	// Requires a set of Requirements indicating what dependency Requirements are needed
	// - if the set is empty, then no requirements are needed
	Requires(context.Context) Requirements
}

// --- Dependency requirement definitiond ---

// Requirements set
type Requirements []Requirement

// UnMatched requirements subset of all Requirements that are not Matched
func (rs Requirements) UnMatched(ctx context.Context) Requirements {
	urs := Requirements{}

	for _, r := range rs {
		if d := r.Matched(ctx); d == nil {
			urs = append(urs, r)
		}
	}

	return urs
}

// Requirement dependency definition from the perspective of the requirer
type Requirement interface {
	// Id unique identifier for the requirement
	Id() string
	// Describe the Requirements for labelling/auditing
	Describe() string
	// Match with a Dependency
	Match(Dependency) error
	// Matched has been Matched. If not matched return nil
	Matched(context.Context) Dependency
}

type OptionalRequirement interface {
	// Optoinal if no error is returned
	Optional() error
}

// CollectRequirements from a set of HasDependencies
func CollectRequirements(ctx context.Context, hds []HasDependencies) Requirements {
	rs := Requirements{}
	for _, hd := range hds {
		if hd == nil {
			continue // stupidity check
		}

		for _, r := range hd.Requires(ctx) {
			slog.Debug("Dependency requirement received", slog.Any("handler", hd), slog.Any("requirement", r))
			rs = append(rs, r)
		}
	}
	return rs
}
