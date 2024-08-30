package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

/**
 * Use this for testing cases where you need to test dependency collection,
 * but in most cases the mock.Component is what you want to use
 */

// RequiresDependencies provide a simple struct that implements dependency.RequiresDependencies.
func RequiresDependencies(rs dependency.Requirements) hasDependencies {
	return hasDependencies{
		rs: rs,
	}
}

type hasDependencies struct {
	rs dependency.Requirements
}

// RequiresDependencies return the requirements.
func (hd hasDependencies) RequiresDependencies(ctx context.Context) dependency.Requirements {
	return hd.rs
}
