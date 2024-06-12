package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	CommandKeyKubernetesConf = "k0s-kubernetes-conf"

	CommandPhaseDiscover       = "K0S Discover"
	CommandPhaseApply          = "K0S Apply"
	CommandPhaseReset          = "K0S Reset"
	CommandPhaseKubernetesConf = "K0S Kubernetes Configuration"
)

func (c *Component) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{ // Requirements that we need
		c.hs,
	}
	ds := dependency.Dependencies{ // Dependencies that our phases typically deliver
		c.k8sd,
		c.k0sd,
	}

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
			&installK0sStep{
				baseStep: bs,
				id:       c.Name(),
			},
			&configureK0sStep{
				baseStep: bs,
				id:       c.Name(),
			},
			&activateK0sStep{
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
			&uninstallK0sStep{
				baseStep: bs,
				id:       c.Name(),
			},
		})
		cmd.Phases.Add(pr)

	case CommandKeyKubernetesConf:
		p := stepped.NewSteppedPhase(CommandPhaseKubernetesConf, rs, dependency.Dependencies{}, []string{dependency.EventKeyActivated})
		p.Steps().Merge(stepped.Steps{
			&discoverStep{
				baseStep: bs,
				id:       c.Name(),
			},
			&kubernetesConfStep{
				baseStep: bs,
				id:       c.Name(),
			},
		})
		cmd.Phases.Add(p)
	}

	return nil
}
