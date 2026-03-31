package mke4

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/creasty/defaults"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeMKE4Component)
}

func DecodeMKE4Component(id string, d func(any) error) (component.Component, error) {
	c := Config{}

	if err := defaults.Set(&c); err != nil {
		return nil, fmt.Errorf("failed to set mke defaults: %w", err)
	}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("failure to decode product '%s' : %w", ComponentType, err)
	}

	if c.Name == "" {
		c.Name = id
	}

	return NewComponent(id, c), nil
}
