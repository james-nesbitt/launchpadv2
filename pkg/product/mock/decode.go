package mock

import (
	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func init() {
	// Register the Mock Products as a decodable product
	product.RegisterDecoder("mock", func(d func(interface{}) error) (component.Component, error) {
		var p prod
		err := d(&p)
		return &p, err
	})
}
