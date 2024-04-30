package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
)

func NewProduct(name string) component.Component {
	return &prod{
		name: name,
	}
}

type prod struct {
	name string
}

func (p *prod) Name() string {
	return p.name
}

func (p prod) Debug() interface{} {
	return struct {
		Name string
	}{
		Name: p.name,
	}
}

func (p *prod) Validate(_ context.Context) error {
	return nil
}
