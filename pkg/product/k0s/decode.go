package k0s

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.RegisterDecoder("k0s", DecodeK0sComponent)
}

func DecodeK0sComponent(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'K0s' : %w", err)
	}

	return NewK0S(id, c), nil
}
