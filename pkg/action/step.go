package action

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Stepped Phases
//
// Sometimes it is convenient to build a phase from a list of functional steps
// instead of as a single function.  Steps allow reusing of functional components
// in different phases.

var (
	ErrSteppedPhaseHasNoSteps = errors.New("stepped phase has no steps")
)

type Steps []Step

// Add a Step to the set of Steps.
func (ss *Steps) Add(sa Step) {
	*ss = append(*ss, sa)
}

// Merge two Steps together.
func (ss *Steps) Merge(ssa Steps) {
	*ss = append(*ss, ssa...)
}

func (ss Steps) Order() (Steps, error) {
	return Steps{}, nil
}

// Step in a stepped Phase.
type Step interface {
	Id() string
	Run(context.Context) error
}

// NewSteppedPhase create an empty phase with steps that can optionally be ordered.
func NewSteppedPhase(id string, orderable bool) steppedPhase {
	return steppedPhase{
		id:        id,
		orderable: orderable,
		steps:     &Steps{},
		Delivers:  Events{},
		Before:    Events{},
		After:     Events{},
	}
}

// steppedPhase a phase which runs a sequence of steps.
type steppedPhase struct {
	id        string
	orderable bool // if true then the steps will be ordered before running
	steps     *Steps
	Delivers  Events
	Before    Events
	After     Events
}

func (sp steppedPhase) Id() string {
	return sp.id
}

func (sp steppedPhase) Validate(ctx context.Context) error {
	if len(*sp.steps) == 0 {
		return ErrSteppedPhaseHasNoSteps
	}

	errs := []error{}
	for _, sps := range *sp.steps {
		spsv, ok := sps.(Validator)
		if !ok {
			continue
		}
		if err := spsv.Validate(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("phase step validation failure: %s", errors.Join(errs...))
	}

	return nil
}

func (sp steppedPhase) Run(ctx context.Context) error {
	var ss Steps = *sp.steps
	if sp.orderable {
		sso, err := ss.Order()
		if err != nil {
			return err
		}
		ss = sso
	}

	for _, s := range ss {
		slog.InfoContext(ctx, fmt.Sprintf("-- STEP: %s", s.Id()), slog.Any("step", s))
		if err := s.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sp steppedPhase) Steps() *Steps {
	return sp.steps
}

func (sp steppedPhase) DeliversEvents(context.Context) Events {
	return sp.Delivers
}

func (sp steppedPhase) RequiresEvents(context.Context) (Events, Events) {
	return sp.Before, sp.After
}
