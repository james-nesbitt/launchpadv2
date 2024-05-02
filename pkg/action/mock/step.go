package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

func NewStep(id string, delivers dependency.Events, requiresBefore dependency.Events, requiresAfter dependency.Events, runError error) step {
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
	resb dependency.Events
	resa dependency.Events
	des  dependency.Events
	rerr error
}

func (s step) Id() string {
	return s.id
}

func (s step) DeliversEvents(_ context.Context) dependency.Events {
	return s.des
}

func (s step) RequiresEvents(_ context.Context) (dependency.Events, dependency.Events) {
	return s.resb, s.resa
}

func (s step) Run(context.Context) error {
	return s.rerr
}
