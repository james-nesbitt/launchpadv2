package k0s

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandKeyKubernetesConf = "k0s-kubernetes-conf"

	CommandPhaseDiscover       = "K0S Discover"
	CommandPhaseApply          = "K0S Apply"
	CommandPhaseReset          = "K0S Reset"
	CommandPhaseKubernetesConf = "K0S Kubernetes Configuration"
)

func (c K0S) BuildCommand(ctx context.Context, cmd *action.Command) error {

	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := action.NewSteppedPhase(CommandPhaseDiscover, false)

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
			&prepaseHostStep{
				id: c.Name(),
			},
			&installK0sStep{
				id: c.Name(),
			},
			&configureK0sStep{
				id: c.Name(),
			},
			&activateK0sStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		pd := action.NewSteppedPhase(CommandPhaseDiscover, false)

		c.dependencyRequiresEvents(ctx, pd.Before, pd.After)

		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := action.NewSteppedPhase(CommandPhaseReset, false)

		c.dependencyDeliversEvents(ctx, dependency.EventKeyDeActivated, pr.Delivers)
		if err := c.dependencyRequiresEvents(ctx, pr.Before, pr.After); err != nil {
			return fmt.Errorf("MKE3 dependency generation error: %w", err)
		}

		pr.Steps().Merge(action.Steps{
			&uninstallK0sStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)

	case CommandKeyKubernetesConf:
		p := action.NewSteppedPhase(CommandPhaseKubernetesConf, false)

		c.dependencyDeliversEvents(ctx, dependency.EventKeyDeActivated, p.Delivers)
		if err := c.dependencyRequiresEvents(ctx, p.Before, p.After); err != nil {
			return fmt.Errorf("MKE3 dependency generation error: %w", err)
		}

		p.Steps().Merge(action.Steps{
			&discoverStep{
				id: c.Name(),
			},
			&kubernetesConfStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(p)
	}

	return nil
}

func (c K0S) dependencyDeliversEvents(ctx context.Context, eventKey string, delivers action.Events) error {
	// collect all events for all delivered dependencies
	ds := map[string]dependency.Dependency{
		"kubernetes": c.k8sd,
		"k0s":        c.k0sd,
	}
	err := action.DependencyDeliversEvents(ctx, []string{eventKey}, ds, delivers)
	if err == nil || errors.Is(err, dependency.ErrNotHandled) {
		return nil
	}

	return fmt.Errorf("failed to build delivered events; %w", err)
}

func (c K0S) dependencyRequiresEvents(ctx context.Context, before, after action.Events) error {
	// collect all events (before and after) for required dependencies
	rs := map[string]dependency.Requirement{
		"controller-hosts": c.chr,
		"worker-hosts":     c.whr,
	}

	if err := action.DependencyRequiresEvents(ctx, rs, before, after); err != nil {
		return fmt.Errorf("%w; failed to build requires events", err)
	}

	return nil
}
