package host

import (
	"context"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/action"
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
	pa := action.NewSteppedPhase(CommandKeyDiscover, false)
	pa.Steps().Add(&discoverStep{
		id: hc.id,
	})

	for _, d := range hc.deps {
		de, ok := d.(action.DeliversEvents)
		if !ok {
			continue
		}

		es := de.DeliversEvents(ctx)

		// we only include dependency activation in discover phases
		if ace, ok := es[dependency.EventKeyActivated]; ok {
			pa.Delivers[ace.Id] = ace
		} else {
			slog.WarnContext(ctx, "host dependency delivers no activation event", slog.Any("dependency", d))
		}
	}

	cmd.Phases.Add(pa)

	return nil
}
