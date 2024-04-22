package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func NewCommandHandler(phases action.Phases, dependencies dependency.Dependencies, events action.Events, buildErr error) cmdHandler {
	return cmdHandler{
		ps:  phases,
		ds:  dependencies,
		es:  events,
		err: buildErr,
	}
}

type cmdHandler struct {
	ds  dependency.Dependencies
	ps  action.Phases
	es  action.Events
	err error
}

func (ch cmdHandler) CommandBuild(_ context.Context, cmd *action.Command) error {
	cmd.Dependencies.Merge(ch.ds)
	cmd.Phases.Merge(ch.ps)
	cmd.Events.Merge(ch.es)

	return ch.err
}
