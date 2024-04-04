package phase

import (
	"context"
	"fmt"
	"strings"
)

type Manager struct {
	phases []Phase
}

// Validate all of the phases.
func (m *Manager) Validate(ctx context.Context) error {
	cctx, cancel := context.WithCancel(ctx)

	errs := []string{}
	for _, p := range m.phases {
		if err := p.Validate(cctx); err != nil {
			errs = append(errs, err.Error())
		}
	}

	cancel()

	if len(errs) > 0 {
		return fmt.Errorf("%w; %s", ErrPhaseValidationFailure, strings.Join(errs, ", "))
	}

	return nil
}

// Run all of the phases.
func (m *Manager) Run(ctx context.Context) error {
	cctx, cancel := context.WithCancel(ctx)

	for _, p := range m.phases {
		if err := p.Run(cctx); err != nil {
			cancel()
			return err
		}
	}

	cancel()

	return nil
}

// Rollback all of the phases.
func (m *Manager) Rollback(ctx context.Context) error {
	cctx, cancel := context.WithCancel(ctx)

	for _, p := range m.phases {
		if err := p.Rollback(cctx); err != nil {
			cancel()
			return err
		}
	}

	cancel()

	return nil
}
