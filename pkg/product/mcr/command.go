package mcr

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandPhaseDiscover = "MCR-Discover"
	CommandPhaseApply    = "MCR-Apply"
	CommandPhaseReset    = "MCR-Reset"
)

func (c MCR) CommandBuild(ctx context.Context, cmd *action.Command) error {
	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyRequiresEvents(ctx, p.Before, p.After)
		c.dependencyDeliversEvents(ctx, dependency.EventKeyActivated, p.Delivers)

		p.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case action.CommandKeyApply:
		p := action.NewSteppedPhase(CommandPhaseApply, false)

		c.dependencyRequiresEvents(ctx, p.Before, p.After)
		c.dependencyDeliversEvents(ctx, dependency.EventKeyActivated, p.Delivers)

		p.Steps().Merge(action.Steps{
			&discoverStep{
				id: c.Name(),
			},
			&prepareMCRNodesStep{
				id: c.Name(),
			},
			&installMCRStep{
				id: c.Name(),
			},
			&swarmActivateStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		pd := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyRequiresEvents(ctx, pd.Before, pd.After)
		c.dependencyDeliversEvents(ctx, dependency.EventKeyActivated, pd.Delivers)

		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := action.NewSteppedPhase(CommandPhaseReset, false)

		c.dependencyRequiresEvents(ctx, pr.Before, pr.After)
		c.dependencyDeliversEvents(ctx, dependency.EventKeyDeActivated, pr.Delivers)

		pr.Steps().Merge(action.Steps{
			&uninstallMCRStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)
	}

	return nil
}

func (c MCR) dependencyDeliversEvents(ctx context.Context, eventKey string, delivers action.Events) error {
	// collect all events for all delivered dependencies
	ds := map[string]dependency.Dependency{
		"dockerhosts": c.dhd,
		"dockerswarm": c.dsd,
	}
	err := action.DependencyDeliversEvents(ctx, []string{eventKey}, ds, delivers)
	if err == nil || errors.Is(err, dependency.ErrNotHandled) {
		return nil
	}

	return fmt.Errorf("failed to build delivered events; %w", err)
}

func (c MCR) dependencyRequiresEvents(ctx context.Context, before, after action.Events) error {
	// collect all events (before and after) for required dependencies
	rs := map[string]dependency.Requirement{
		"manager-hosts": c.mrhr,
		"worker-hosts":  c.wrhr,
	}

	if err := action.DependencyRequiresEvents(ctx, rs, before, after); err != nil {
		return fmt.Errorf("%w; failed to build requires events", err)
	}

	return nil
}
