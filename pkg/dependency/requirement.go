package dependency

import (
	"context"
)

// --- Dependency requirement definitiond ---

// Requirements set
type Requirements []Requirement

// UnMet dependency requirements set
func (rs Requirements) UnMet(ctx context.Context) Requirements {
	urs := Requirements{}

	for _, r := range rs {
		if err := r.Validate(ctx); err != nil {
			urs = append(urs, r)
		}
	}

	return urs
}

// Requirement dependency definition from the perspective of the requirer
type Requirement interface {
	// Describe the requirements for labelling/auditing
	Describe() string
	// Validate that the requirement can be met
	//   Validation does not meet the dependency, it only confirms that the dependency can be met
	//   when it is needed.  Each requirement will have its own interface for individual types of
	//   requirements.
	Validate(context.Context) error
}

// FullfillingRequirement a requirement that needs to be fullfilled
type FullfillingRequirement interface {
	// Fullfilled the requirement fullfillment is performed
	//   This indicates that the requirement deliverebles are now available, indicating that any
	//   requirer of the dependency can now utilize
	Fullfilled(context.Context) bool
}
