package msr2

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/creasty/defaults"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeComponent)
}

func DecodeComponent(id string, d func(any) error) (component.Component, error) {
	c := Config{}

	if err := defaults.Set(&c); err != nil {
		return nil, fmt.Errorf("failed to set msr2 defaults: %w", err)
	}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponent(id, c), nil
}
