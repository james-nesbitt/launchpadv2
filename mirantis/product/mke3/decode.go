package mke3

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

	defaults.Set(&c)

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("failure to unmarshal product 'MKE3' : %w", err)
	}

	return NewComponent(id, c), nil
}
