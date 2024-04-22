package msr3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
)

const (
	CommandPhaseDiscover = "MSR3 Discover"
	CommandPhaseApply    = "MSR3 Apply"
	CommandPhaseReset    = "MSR3 Reset"
)

func (c MSR3) BuildCommand(ctx context.Context, cmd *action.Command) error {

	switch cmd.Key {
	case action.CommandKeyDiscover:
		p := action.NewSteppedPhase(action.CommandPhaseDiscover, false)
		p.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case action.CommandKeyApply:
		p := action.NewSteppedPhase(CommandPhaseApply, false)
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
		})
		cmd.Phases.Add(p)

	case action.CommandKeyReset:
		pd := action.NewSteppedPhase(action.CommandPhaseDiscover, false)
		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := action.NewSteppedPhase(CommandPhaseReset, false)
		pr.Steps().Merge(action.Steps{
			&uninstallMSRStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)

	}

	return nil
}
