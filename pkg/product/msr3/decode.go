package msr3

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.RegisterDecoder("msr3", DecodeMSR3Component)
}

func DecodeMSR3Component(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MSR3' : %w", err)
	}

	return NewMSR3(id, c), nil
}
