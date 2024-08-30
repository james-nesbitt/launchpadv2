package dependency

import (
	"context"
	"errors"
)

var (
	ErrRequirementNotCreated = errors.New("requirement not created") // A requirement exists but was never initialized (is nil)
	ErrRequirementNotMatched = errors.New("requirement not matched") // a Requirement exists but has no Dependency Matched (is nil)
)

// HasDependencies can determine its dependencies.
type RequiresDependencies interface {
	// Requires a set of Requirements indicating what dependency Requirements are needed
	// - if the set is empty, then no requirements are needed
	RequiresDependencies(context.Context) Requirements
}

// --- Dependency requirement definitiond ---

// Requirements set.
type Requirements []Requirement

// NewRequirements a new Requirements, optionally pre-populated.
func NewRequirements(rsa ...Requirement) Requirements {
	rs := Requirements{}
	if len(rsa) > 0 {
		rs = append(rs, rsa...)
	}
	return rs
}

// UnMatched requirements subset of all Requirements that are not Matched.
func (rs Requirements) UnMatched(ctx context.Context) Requirements {
	urs := Requirements{}

	for _, r := range rs {
		if d := r.Matched(ctx); d == nil {
			urs = append(urs, r)
		}
	}

	return urs
}
