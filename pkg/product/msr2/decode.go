package msr2

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.RegisterDecoder("msr2", DecodeMSR2Component)
}

func DecodeMSR2Component(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MSR2' : %w", err)
	}

	return NewMSR2(id, c), nil
}
