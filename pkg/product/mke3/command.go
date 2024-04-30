package mke3

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandPhaseDiscover = "MKE3-Discover"
	CommandPhaseApply    = "MKE3-Apply"
	CommandPhaseReset    = "MKE3-Reset"
)

func (c MKE3) CommandBuild(ctx context.Context, cmd *action.Command) error {
	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyDeliversEvents(ctx, dependency.EventKeyActivated, p.Delivers)
		c.dependencyRequiresEvents(ctx, p.Before, p.After)

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
			return fmt.Errorf("MKE3 dependency generation error: %w", err)
		}

		p.Steps().Merge(action.Steps{
			&discoverStep{
				id: c.Name(),
			},
			&prepareNodesStep{
				id: c.Name(),
			},
			&installMKEStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		pd := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyDeliversEvents(ctx, dependency.EventKeyActivated, pd.Delivers)
		c.dependencyRequiresEvents(ctx, pd.Before, pd.After)

		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := action.NewSteppedPhase(CommandPhaseReset, false)

		c.dependencyRequiresEvents(ctx, pr.Before, pr.After)
		if err := c.dependencyDeliversEvents(ctx, dependency.EventKeyDeActivated, pr.Delivers); err != nil {
			return fmt.Errorf("MKE3 dependency generation error: %w", err)
		}

		pr.Steps().Merge(action.Steps{
			&uninstallMKEStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)
	}

	return nil
}

func (c MKE3) dependencyDeliversEvents(ctx context.Context, eventKey string, delivers action.Events) error {
	// collect all events for all delivered dependencies
	ds := map[string]dependency.Dependency{
		"kubernetes": c.k8sd,
		"MKE3":       c.mked,
	}
	err := action.DependencyDeliversEvents(ctx, []string{eventKey}, ds, delivers)
	if err == nil || errors.Is(err, dependency.ErrNotHandled) {
		return nil
	}

	return fmt.Errorf("failed to build delivered events; %w", err)
}

func (c MKE3) dependencyRequiresEvents(ctx context.Context, before, after action.Events) error {
	// collect all events (before and after) for required dependencies
	rs := map[string]dependency.Requirement{
		"dockerhosts": c.dhr,
	}

	if err := action.DependencyRequiresEvents(ctx, rs, before, after); err != nil {
		return fmt.Errorf("%w; failed to build requires events", err)
	}

	return nil
}
