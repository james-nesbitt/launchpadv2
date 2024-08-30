package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// Custom Requirement types

// StaticRequirement provides a simple requirement with a dependency already injected, that does not allow new matches.
func StaticRequirement(id, desc string, d dependency.Dependency) *reqDepMatcher {
	return &reqDepMatcher{
		id:   id,
		desc: desc,
		matcher: func(_ dependency.Dependency) error {
			return dependency.ErrDependencyWrongType
		},
		d: d,
	}
}

// SimpleRequirement provides a simple requirement that takes any dependency.
func SimpleRequirement(id, desc string) *reqDepMatcher {
	return &reqDepMatcher{
		id:   id,
		desc: desc,
	}
}

// MatchingRequirement provide a requirement that accepts dependencies based on a provided mathching function.
func MatchingRequirement(id, desc string, matcher func(dependency.Dependency) error) *reqDepMatcher {
	return &reqDepMatcher{
		id:      id,
		desc:    desc,
		matcher: matcher,
	}
}

type reqDepMatcher struct {
	id      string
	desc    string
	matcher func(dependency.Dependency) error

	d dependency.Dependency
}

// Id unique identifier for the requirement.
func (r reqDepMatcher) Id() string {
	return r.id
}

// Describe the Requirements for labelling/auditing.
func (r reqDepMatcher) Describe() string {
	return r.desc
}

// Match with a Dependency.
func (r *reqDepMatcher) Match(d dependency.Dependency) error {
	if r.matcher != nil {
		if err := r.matcher(d); err != nil {
			return err
		}
	}

	r.d = d
	return nil
}

// Matched has been Matched. If not matched return nil.
func (r reqDepMatcher) Matched(context.Context) dependency.Dependency {
	return r.d
}
