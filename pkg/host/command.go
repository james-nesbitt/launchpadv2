package host

import (
	"context"
	"fmt"
	"log/slog"

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
	bs := baseStep{c: hc}

	if len(hc.deps) > 0 { // only discover if something else needs us
		p := stepped.NewSteppedPhase(CommandKeyDiscover, dependency.Requirements{}, hc.deps, []string{dependency.EventKeyActivated})

		for _, d := range hc.deps {
			if d == nil {
				slog.WarnContext(ctx, "HostsComponent has a nil Dependency", slog.Any("component", hc))
				continue // sanity check
			}

			p.Steps().Add(&discoverStep{
				baseStep: bs,
				id:       fmt.Sprintf("%s", d.Id()),
				d:        d,
			})
		}

		cmd.Phases.Add(p)
	}

	return nil
}
