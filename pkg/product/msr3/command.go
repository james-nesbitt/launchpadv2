package msr3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	CommandPhaseDiscover = "MSR3 Discover"
	CommandPhaseApply    = "MSR3 Apply"
	CommandPhaseReset    = "MSR3 Reset"
)

func (c Component) BuildCommand(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need
		c.k8sr,
	}
	ds := dependency.Dependencies{}

	switch cmd.Key {
	case project.CommandKeyDiscover:
		p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case project.CommandKeyApply:
		p := stepped.NewSteppedPhase(CommandPhaseApply, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
			&discoverStep{
				id: c.Name(),
			},
			&installMSRStep{
				id: c.Name(),
			},
			&configureMSRStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(p)

	case project.CommandKeyReset:
		pd := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := stepped.NewSteppedPhase(CommandPhaseReset, rs, ds, []string{dependency.EventKeyDeActivated})
		pr.Steps().Merge(stepped.Steps{
			&uninstallMSRStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)

	}

	return nil
}
