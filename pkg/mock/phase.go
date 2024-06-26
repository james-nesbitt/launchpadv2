package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

func NewPhase(id string, run, validate func(context.Context) error, delivers, before, after dependency.Events) phase {
	return phase{
		id:       id,
		run:      run,
		validate: validate,
		delivers: delivers,
		before:   before,
		after:    after,
	}
}

type phase struct {
	id       string
	run      func(context.Context) error
	validate func(context.Context) error
	delivers dependency.Events
	before   dependency.Events
	after    dependency.Events
}

func (p phase) Id() string {
	return p.id
}

func (p phase) Validate(ctx context.Context) error {
	return p.validate(ctx)
}

func (p phase) Run(ctx context.Context) error {
	return p.run(ctx)
}

func (p phase) DeliversEvents(_ context.Context) dependency.Events {
	return p.delivers
}

func (p phase) RequiresEvents(_ context.Context) (dependency.Events, dependency.Events) {
	return p.before, p.after
}
