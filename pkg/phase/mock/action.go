package mock

import (
	"context"
)

func NewAction(label string, validate, run func(context.Context) error) action {
	return action{
		label:    label,
		validate: validate,
		run:      run,
	}
}

type action struct {
	label string

	validate func(context.Context) error
	run      func(context.Context) error
}

func (a action) Label() string {
	return a.label
}

func (a action) Validate(ctx context.Context) error {
	return a.validate(ctx)
}

func (a action) Run(ctx context.Context) error {
	return a.run(ctx)
}

func NewReversableAction(label string, validate, run, rollback func(context.Context) error) rollbackAction {
	return rollbackAction{
		action: action{
			label:    label,
			validate: validate,
			run:      run,
		},
		rollback: rollback,
	}
}

type rollbackAction struct {
	action

	rollback func(context.Context) error
}

func (a rollbackAction) Rollback(ctx context.Context) error {
	return a.rollback(ctx)
}
