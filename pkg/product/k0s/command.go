package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	CommandKeyKubernetesConf = "k0s-kubernetes-conf"

	CommandPhaseDiscover       = "K0S Discover"
	CommandPhaseApply          = "K0S Apply"
	CommandPhaseReset          = "K0S Reset"
	CommandPhaseKubernetesConf = "K0S Kubernetes Configuration"
)

func (c K0S) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need
		c.chr,
		c.whr,
	}
	ds := dependency.Dependencies{ // Dependencies that our phases typically deliver
		c.k8sd,
		c.k0sd,
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
		pd := stepped.NewSteppedPhase(CommandPhaseApply, rs, ds, []string{dependency.EventKeyActivated})
		pd.Steps().Add(
			&discoverStep{
				id: c.Name(),
			},
		)
		cmd.Phases.Add(pd)

		pr := stepped.NewSteppedPhase(CommandPhaseReset, rs, ds, []string{dependency.EventKeyDeActivated})
		pr.Steps().Merge(stepped.Steps{
			&uninstallK0sStep{
				id: c.Name(),
			},
		})
		cmd.Phases.Add(pr)

	case CommandKeyKubernetesConf:
		p := stepped.NewSteppedPhase(CommandPhaseKubernetesConf, rs, dependency.Dependencies{}, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
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
