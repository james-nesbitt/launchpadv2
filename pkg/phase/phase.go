package phase

import (
	"context"
	"errors"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type Phase interface {
	// Name the stage for logging and auditing
	Name() string
	// Validate the phase and all of its actions without actually interfering with any systems
	Validate(context.Context) error
	// Run all of the actions for the phase
	Run(context.Context) error
	// Rollback all of the actions for the phase (May do nothing) - maybe this needs it own interphase?
	Rollback(context.Context) error
}

var (
	ErrPhaseValidationFailure = errors.New("phase validation failed")
	ErrPhaseRunFailure        = errors.New("phase execution failed")
	ErrPhaseRollbackFailure   = errors.New("phase rollback failed")
)

type ProvidesActions interface {
	Name() string

	ProvideActions(ctx context.Context, phaseId string, deps dependency.Dependencies) (Actions, error)
}

type ProvidesPhases interface {
	SuggestPhases(command string) error
}
