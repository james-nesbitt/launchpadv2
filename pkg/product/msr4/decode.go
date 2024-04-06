package msr4

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.RegisterDecoder("msr4", DecodeMSR4Component)
}

func DecodeMSR4Component(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MSR4' : %w", err)
	}

	return NewMSR4(id, c), nil
}
