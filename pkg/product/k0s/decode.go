package k0s

import (
	"fmt"

	"github.com/k0sproject/dig"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	product.RegisterDecoder(ComponentType, DecodeComponent)
}

// DecodeComponent decode a new component from an unmarshall decoder
func DecodeComponent(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{
		K0sConfig: dig.Mapping{},
	}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponent(id, c), nil
}
