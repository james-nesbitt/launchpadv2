package product_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/product/mock"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/phase"
	"github.com/Mirantis/launchpad/pkg/product"
)

func Test_Decode(t *testing.T) {
	product.RegisterDecoder("dummy", func(func(interface{}) error) (component.Component, error) {
		p := mock.NewProduct(
			"test",
			phase.Actions{},
		)
		return p, nil
	})

	c, err := product.DecodeKnownProduct("dummy", nil)
	if err != nil {
		t.Errorf("Unexpected decode error: %s", err.Error())
	} else if c.Name() != "test" {
		t.Errorf("Product wasn't what we expected: %+v", c)
	}
}
