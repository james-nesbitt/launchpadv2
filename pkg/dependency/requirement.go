package dependency

import (
	"context"
	"log/slog"
)

// Requirement dependency definition from the perspective of the requirer.
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

// CollectRequirements from a set of HasDependencies.
func CollectRequirements(ctx context.Context, hds []RequiresDependencies) Requirements {
	rs := Requirements{}
	for _, hd := range hds {
		if hd == nil {
			continue // stupidity check
		}

		for _, r := range hd.RequiresDependencies(ctx) {
			slog.Debug("Dependency requirement received", slog.Any("handler", hd), slog.Any("requirement", r))
			rs = append(rs, r)
		}
	}
	return rs
}
