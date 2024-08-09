package msr2

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeComponent)
}

func DecodeComponent(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponent(id, c), nil
}
