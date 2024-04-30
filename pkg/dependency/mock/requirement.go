package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

func HasDependencies(rs dependency.Requirements) dependency.HasDependencies {
	return hasDependencies{
		rs: rs,
	}
}

type hasDependencies struct {
	rs dependency.Requirements
}

func (mhc hasDependencies) Requires(ctx context.Context) dependency.Requirements {
	return mhc.rs
}

func Requirement(id, description string, validate error) dependency.Requirement {
	return &req{
		id:          id,
		description: description,
		err:         validate,
	}
}

type req struct {
	id          string
	description string
	err         error
	d           dependency.Dependency
}

func (mr req) Id() string {
	return mr.id
}

func (mr req) Describe() string {
	return mr.description
}

func (mr req) Validate(_ context.Context) error {
	return mr.err
}

func (mr *req) Match(d dependency.Dependency) error {
	mr.d = d
	return nil
}

func (mr req) Matched(_ context.Context) dependency.Dependency {
	return mr.d
}
