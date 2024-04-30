package product_test

import (
	"strings"
	"testing"

	"github.com/Mirantis/launchpad/pkg/product/mock"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product"
)

func Test_Decode(t *testing.T) {
	product.RegisterDecoder("dummy", func(id string, _ func(interface{}) error) (component.Component, error) {
		p := mock.NewProduct(id)
		return p, nil
	})

	c, err := product.DecodeKnownProduct("dummy", "dummy", nil)
	if err != nil {
		t.Errorf("Unexpected decode error: %s", err.Error())
	} else if !strings.Contains(c.Name(), "dummy") {
		t.Errorf("Product wasn't what we expected: %+v", c)
	}
}
