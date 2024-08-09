package mke3

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
		return nil, fmt.Errorf("Failure to unmarshal product 'MKE3' : %w", err)
	}

	return NewComponent(id, c), nil
}
