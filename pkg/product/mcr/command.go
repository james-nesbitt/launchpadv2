package mcr

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandPhaseDiscover = "MCR-Discover"
	CommandPhaseApply    = "MCR-Apply"
	CommandPhaseReset    = "MCR-Reset"
)

func (c *MCR) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need
		c.hr,
	}
	ds := dependency.Dependencies{ // Dependencies that we deliver
		c.dhd,
		c.dsd,
	}

	bs := baseStep{
		c: c,
	}

	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Add(
			&discoverStep{
				id:          c.Name(),
				baseStep:    bs,
				failOnError: true,
			},
		)
		cmd.Phases.Add(p)

	case action.CommandKeyApply:
		p := stepped.NewSteppedPhase(CommandPhaseApply, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
			&installMCRStep{
				id:       c.Name(),
				baseStep: bs,
			},
			&swarmActivateStep{
				id:       c.Name(),
				baseStep: bs,
			},
			&discoverStep{
				id:          fmt.Sprintf("%s-after", c.Name()),
				baseStep:    bs,
				failOnError: true,
			},
		})
		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		if len(ds) > 0 { // only discover if something else needs us - so that we can meet any dependencies

			pd := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
			pd.Steps().Add(
				&discoverStep{
					id:       c.Name(),
					baseStep: bs,
				},
			)
			cmd.Phases.Add(pd)
		}

		pr := stepped.NewSteppedPhase(CommandPhaseReset, rs, ds, []string{dependency.EventKeyDeActivated})
		pr.Steps().Merge(stepped.Steps{
			&uninstallMCRStep{
				id:       c.Name(),
				baseStep: bs,
			},
			&discoverStep{
				id:       fmt.Sprintf("%s-after", c.Name()),
				baseStep: bs,
			},
		})
		cmd.Phases.Add(pr)
	}

	return nil
}
