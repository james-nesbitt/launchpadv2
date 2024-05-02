package stepped

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

var (
	ErrSteppedPhaseHasNoSteps = errors.New("stepped phase has no steps")
)

// NewSteppedPhase create an empty phase with steps that can optionally be ordered.
func NewSteppedPhase(id string, requires dependency.Requirements, delivers dependency.Dependencies, deliverEventKeys []string) SteppedPhase {
	return SteppedPhase{
		id:    id,
		steps: &Steps{},
		reqs:  requires,
		deps:  delivers,
		deks:  deliverEventKeys,
	}
}

// SteppedPhase a phase which runs a sequence of steps.
type SteppedPhase struct {
	id    string
	steps *Steps
	reqs  dependency.Requirements
	deps  dependency.Dependencies
	deks  []string
}

func (sp SteppedPhase) Id() string {
	return sp.id
}

func (sp SteppedPhase) Validate(ctx context.Context) error {
	errs := []error{}

	if _, err := sp.deliveredEventsCheck(ctx); err != nil {
		errs = append(errs, err)
	}
	if _, _, err := sp.requiredEventsCheck(ctx); err != nil {
		errs = append(errs, err)
	}

	// Validate all of the steps
	if len(*sp.steps) == 0 {
		return ErrSteppedPhaseHasNoSteps
	}
	for _, sps := range *sp.steps {
		spsv, ok := sps.(action.Validator)
		if !ok {
			continue
		}
		if err := spsv.Validate(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("phase '%s' validation failure: %s", sp.Id(), errors.Join(errs...))
	}

	return nil
}

func (sp SteppedPhase) Run(ctx context.Context) error {
	var ss Steps = *sp.steps

	for _, s := range ss {
		slog.InfoContext(ctx, fmt.Sprintf("-- STEP: %s", s.Id()), slog.Any("step", s))
		if err := s.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (sp SteppedPhase) Steps() *Steps {
	return sp.steps
}
