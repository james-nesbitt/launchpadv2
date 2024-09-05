package mke4

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	CommandPhaseDiscover = "MKE4-Discover"
	CommandPhaseApply    = "MKE4-Apply"
	CommandPhaseReset    = "MKE4-Reset"
)

func (c *Component) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need to operate
		c.k8sr, // MKE4 can install on any kubernetes context
	}
	ds := dependency.Dependencies{} // Dependencies that we deliver

	bs := baseStep{
		c: c,
	}

	switch cmd.Key {
	case project.CommandKeyDiscover:
		p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Add(
			&discoverStep{
				baseStep: bs,
				id:       c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case project.CommandKeyApply:
		p := stepped.NewSteppedPhase(CommandPhaseApply, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
			&discoverStep{
				baseStep: bs,
				id:       c.Name(),
			},
			&installOperatorStep{
				baseStep: bs,
				id:       c.Name(),
			},
			&activateStep{
				baseStep: bs,
				id:       c.Name(),
			},
		})
		cmd.Phases.Add(p)

	case project.CommandKeyReset:
		if len(ds) > 0 { // only discover if something else needs us
			pd := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
			pd.Steps().Add(
				&discoverStep{
					baseStep: bs,
					id:       c.Name(),
				},
			)
			cmd.Phases.Add(pd)
		}

		pr := stepped.NewSteppedPhase(CommandPhaseReset, rs, ds, []string{dependency.EventKeyDeActivated})
		pr.Steps().Merge(stepped.Steps{
			&deactivateStep{
				baseStep: bs,
				id:       c.Name(),
			},
			&uninstallOperatorStep{
				baseStep: bs,
				id:       c.Name(),
			},
		})
		cmd.Phases.Add(pr)
	}

	return nil
}
