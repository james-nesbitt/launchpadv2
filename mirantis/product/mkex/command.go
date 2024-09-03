package mkex

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	CommandPhaseDiscover = "K0S Discover"
	CommandPhaseApply    = "K0S Apply"
)

func (c *Component) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need
		c.hsr,
	}
	ds := dependency.Dependencies{ // Dependencies that we deliver
		c.dhd,
	}

	bs := baseStep{
		c: c,
	}

	switch cmd.Key {
	case project.CommandKeyDiscover:
		p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Add(
			&discoverStep{
				id:       c.Name(),
				baseStep: bs,
			},
		)
		cmd.Phases.Add(p)

	case project.CommandKeyApply:
		p := stepped.NewSteppedPhase(CommandPhaseApply, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
			&swarmActivateStep{
				id:       c.Name(),
				baseStep: bs,
			},
			&discoverStep{
				id:       fmt.Sprintf("%s-discover-after-%s", c.Name(), cmd.Key),
				baseStep: bs,
			},
		})
		cmd.Phases.Add(p)

	}

	return nil
}
