package mock

import (
	"context"
)

// StaticDependency return a simple dependency with static values, including validate and met responses.
func StaticDependency(id, desc string, validate error, met error) *dep {
	return &dep{
		id: id,
		d:  desc,
		v:  func(_ context.Context) error { return validate },
		m:  func(_ context.Context) error { return met },
	}
}

// SimpleDependency return a simple dependency that uses passed handlers for validation and met.
func SimpleDependency(id, desc string, validate, met func(context.Context) error) *dep {
	return &dep{
		id: id,
		d:  desc,
		v:  validate,
		m:  met,
	}
}

// dep a simple mock dependency that can be programmed how to respond to validation / met requests.
type dep struct {
	id string
	d  string                      // describe
	v  func(context.Context) error // validate handler
	m  func(context.Context) error // met handler
}

func (md dep) Id() string {
	return md.id
}

func (md dep) Describe() string {
	return md.d
}

func (md dep) Validate(ctx context.Context) error {
	return md.v(ctx)
}

func (md dep) Met(ctx context.Context) error {
	return md.m(ctx)
}
