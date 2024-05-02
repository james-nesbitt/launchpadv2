package host

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

var (
	CommandKeyDiscover   = "hosts discover"
	CommandKeyDisconnect = "hosts disconnect"
)

// CommandBuild build HostComponent actions into any command
//
//	For all commands, we effectively do discovery for the hosts
func (hc *HostsComponent) CommandBuild(ctx context.Context, cmd *action.Command) error {
	p := stepped.NewSteppedPhase(CommandKeyDiscover, dependency.Requirements{}, hc.deps, []string{dependency.EventKeyActivated})

	p.Steps().Add(&discoverStep{
		id: hc.id,
	})

	cmd.Phases.Add(p)

	return nil
}
