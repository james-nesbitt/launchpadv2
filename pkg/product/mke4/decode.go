package mke4

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.RegisterDecoder("mke4", DecodeMKE4Component)
}
func DecodeMKE4Component(d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MKE4' : %w", err)
	}

	return NewMKE4(c), nil
}
