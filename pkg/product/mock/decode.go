package mock

import (
	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	// Register the Mock Products as a decodable product
	product.RegisterDecoder("mock", func(id string, d func(interface{}) error) (component.Component, error) {
		var p prod
		err := d(&p)
		p.name = id

		return &p, err
	})
}
