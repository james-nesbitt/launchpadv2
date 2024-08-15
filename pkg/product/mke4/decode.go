package mke4

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/creasty/defaults"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeMKE4Component)
}
func DecodeMKE4Component(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	defaults.Set(&c)

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponent(id, c), nil
}
