package product_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/phase"
	"github.com/Mirantis/launchpad/pkg/product"
)

type DummyProduct struct {
	name    string
	actions phase.Actions
}

func (dp *DummyProduct) Name() string {
	return dp.name
}

func (dp *DummyProduct) Actions(context.Context, string) (phase.Actions, error) {
	return dp.actions, nil
}

func Test_Decode(t *testing.T) {
	product.ProductDecoders["dummy"] = func(product.Decoder) (component.Component, error) {
		return &DummyProduct{
			name:    "test",
			actions: phase.Actions{},
		}, nil
	}

	c, err := product.DecodeKnownProduct("dummy", nil)
	if err != nil {
		t.Errorf("Unexpected decode error: %s", err.Error())
	} else if c.Name() != "test" {
		t.Errorf("Product wasn't what we expected: %+v", c)
	}
}
