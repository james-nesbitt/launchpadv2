package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/phase"
)

func NewProduct(name string, actions phase.Actions) component.Component {
	return &prod{
		name:    name,
		actions: actions,
	}
}

type prod struct {
	name    string
	actions phase.Actions
}

func (p *prod) Name() string {
	return p.name
}

func (p *prod) Actions(context.Context, string) (phase.Actions, error) {
	return p.actions, nil
}
