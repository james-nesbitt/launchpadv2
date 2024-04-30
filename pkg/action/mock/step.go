package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action"
)

func NewStep(id string, delivers action.Events, requiresBefore action.Events, requiresAfter action.Events, runError error) step {
	return step{
		id:   id,
		des:  delivers,
		resb: requiresBefore,
		resa: requiresAfter,
		rerr: runError,
	}
}

type step struct {
	id   string
	resb action.Events
	resa action.Events
	des  action.Events
	rerr error
}

func (s step) Id() string {
	return s.id
}

func (s step) DeliversEvents(_ context.Context) action.Events {
	return s.des
}

func (s step) RequiresEvents(_ context.Context) (action.Events, action.Events) {
	return s.resb, s.resa
}

func (s step) Run(context.Context) error {
	return s.rerr
}
