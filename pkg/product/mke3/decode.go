package mke3

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.ProductDecoders["mke3"] = DecodeMKE3Component
}

func DecodeMKE3Component(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MKE3' : %w", err)
	}

	return NewMKE3(id, c), nil
}
