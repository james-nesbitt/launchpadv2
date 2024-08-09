package component_test

import (
	"strings"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_Decode(t *testing.T) {
	component.RegisterDecoder("dummy", func(id string, _ func(interface{}) error) (component.Component, error) {
		p := mock.NewComponent(id, nil, nil)
		return p, nil
	})

	c, err := component.DecodeComponent("dummy", "dummy", nil)
	if err != nil {
		t.Errorf("Unexpected decode error: %s", err.Error())
	} else if !strings.Contains(c.Name(), "dummy") {
		t.Errorf("Product wasn't what we expected: %+v", c)
	}
}
