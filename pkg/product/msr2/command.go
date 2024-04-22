package msr2

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandPhaseDiscover = "MSR2-Discover"
	CommandPhaseApply    = "MSR2-Apply"
	CommandPhaseReset    = "MSR2-Reset"
)

func (c MSR2) CommandBuild(ctx context.Context, cmd *action.Command) error {

	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyRequiresEvents(ctx, p.Before, p.After)
		if err := c.dependencyRequiresEvents(ctx, p.Before, p.After); err != nil {
			return fmt.Errorf("MSR2 dependency generation error: %w", err)
		}

		p.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case action.CommandKeyApply:
		p := action.NewSteppedPhase(CommandPhaseApply, false)

		c.dependencyDeliversEvents(ctx, dependency.EventKeyActivated, p.Delivers)
		if err := c.dependencyRequiresEvents(ctx, p.Before, p.After); err != nil {
			return fmt.Errorf("MSR2 dependency generation error: %w", err)
		}

		p.Steps().Merge(action.Steps{
			&discoverStep{
				id: c.Name(),
			},
			&installMSRStep{
				id: c.Name(),
			},
			&configureMSRStep{
				id: c.Name(),
			},
			&activateMSRStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		pd := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyRequiresEvents(ctx, pd.Before, pd.After)
		if err := c.dependencyRequiresEvents(ctx, pd.Before, pd.After); err != nil {
			return fmt.Errorf("MSR2 dependency generation error: %w", err)
		}

		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := action.NewSteppedPhase(CommandPhaseReset, false)

		c.dependencyDeliversEvents(ctx, dependency.EventKeyDeActivated, pr.Delivers)
		if err := c.dependencyRequiresEvents(ctx, pr.Before, pr.After); err != nil {
			return fmt.Errorf("MSR2 dependency generation error: %w", err)
		}

		pr.Steps().Merge(action.Steps{
			&uninstallMSRStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)
	}

	return nil
}

func (c MSR2) dependencyDeliversEvents(ctx context.Context, eventKey string, delivers action.Events) error {
	// collect all events for all delivered dependencies
	ds := map[string]dependency.Dependency{
		"msr": c.msr2,
	}
	err := action.DependencyDeliversEvents(ctx, []string{eventKey}, ds, delivers)
	if err == nil || errors.Is(err, dependency.ErrNotHandled) {
		return nil
	}

	return fmt.Errorf("failed to build delivered events; %w", err)
}

func (c MSR2) dependencyRequiresEvents(ctx context.Context, before, after action.Events) error {
	// collect all events (before and after) for required dependencies
	rs := map[string]dependency.Requirement{
		"mke3":        c.mke3r,
		"dockerhosts": c.dhr,
	}

	if err := action.DependencyRequiresEvents(ctx, rs, before, after); err != nil {
		return fmt.Errorf("%w; failed to build requires events", err)
	}

	return nil
}
