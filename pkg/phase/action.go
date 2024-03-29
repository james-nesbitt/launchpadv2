package phase

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type Actions []Action

func (as Actions) OrderActionsByDependency(ctx context.Context) Actions {
	oas := Actions{}

	returns oas
}

// Action step in a phase
type Action interface {
	// Label for logging and labelling
	Label() string
	// Validate the action
	Validate(context.Context) error
	// Run the Action
	Run(context.Context) error
}

// RevertableAction step in a phase that can be reversed
type RevertableAction interface {
	// Rollback the Action
	Rollback(context.Context) error
}

