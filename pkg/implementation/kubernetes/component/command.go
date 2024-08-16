package kubernetes

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	CommandKeyKubernetesConf = "Kubernetes Configuration"

	CommandPhaseDiscover       = "Kubernetes Discover"
	CommandPhaseApply          = "Kubernetes Apply"
	CommandPhaseReset          = "Kubernetes Reset"
	CommandPhaseKubernetesConf = "Kubernetes Configuration"
)

func (c *Component) CommandBuild(ctx context.Context, cmd *action.Command) error {
	rs := dependency.Requirements{} // Requirements that we need
	ds := dependency.Dependencies{  // Dependencies that our phases typically deliver
		c.k8sd,
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
		p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
		p.Steps().Add(
			&discoverStep{
				baseStep: bs,
				id:       c.Name(),
			},
		)
		cmd.Phases.Add(p)

	case project.CommandKeyReset:
		if c.k8sd != nil { // only discover if someone has asked for
			p := stepped.NewSteppedPhase(CommandPhaseDiscover, rs, ds, []string{dependency.EventKeyActivated})
			p.Steps().Add(
				&discoverStep{
					baseStep: bs,
					id:       c.Name(),
				},
			)
			cmd.Phases.Add(p)
		}

	case CommandKeyKubernetesConf:
		p := stepped.NewSteppedPhase(CommandPhaseKubernetesConf, rs, ds, []string{dependency.EventKeyActivated})
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
