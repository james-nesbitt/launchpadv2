package mke3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandPhaseDiscover = "MKE3-Discover"
	CommandPhaseApply    = "MKE3-Apply"
	CommandPhaseReset    = "MKE3-Reset"
)

func (c *MKE3) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need
		c.dhr,
	}
	ds := dependency.Dependencies{ // Dependencies that we deliver
		c.k8sd,
		c.mked,
	}

	bs := baseStep{
		c: c,
	}

	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case action.CommandKeyApply:
		p := stepped.NewSteppedPhase(CommandPhaseApply, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
			&discoverStep{
				id: c.Name(),
			},
			&prepareNodesStep{
				id: c.Name(),
			},
		})

		if c.config.imagePrepullRequired() {
			p.Steps().Add(
				&prePullImagesStep{
					baseStep: bs,
					id:       c.Name(),
				},
			)
		}

		p.Steps().Add(
			&runMKEBootstrapper{
				baseStep: bs,
				id:       c.Name(),
			},
		)

		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		if len(ds) > 0 { // only discover if something else needs us
			pd := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
			pd.Steps().Add(
				&discoverStep{
					id: c.Name(),
				},
			)
			cmd.Phases.Add(pd)
		}

		pr := stepped.NewSteppedPhase(CommandPhaseReset, rs, ds, []string{dependency.EventKeyDeActivated})
		pr.Steps().Merge(stepped.Steps{
			&uninstallMKEStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)
	}

	return nil
}
