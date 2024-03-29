package phase

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	ErrPhaseValidationFailure = errors.New("Phase validation failed")
	ErrPhaseRunFailure        = errors.New("Phase execution failed")
)

type ProvidesActions interface {
	Name() string
	ProvideActions(context.Context) Actions
}

type StandardPhase struct {
	name string
	id   string
	pas  []ProvidesActions
}

func (sp StandardPhase) Name() string {
	return sp.name
}

func (sp StandardPhase) Validate(ctx context.Context) error {
	errs := []string{}

	for _, pa := range sp.pas {
		for _, a := range pa.ProvideActions(ctx) {
			if err := a.Validate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("%s:%s %s", pa.Name(), a.Label(), err.Error()).Error())
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w; phase action validation errors: %s", ErrPhaseValidationFailure, strings.Join(errs, "\n"))
	}
	return nil
}

func (sp StandardPhase) Run(ctx context.Context) error {
	// executed actions, in case we need to rollback
	return nil
}

func (sp StandardPhase) Rollback(ctx context.Context) error {
	// executed actions, in case we need to rollback
	return nil
}
