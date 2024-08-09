package mock

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
)

func init() {
	component.RegisterDecoder("mock", DecodeMockComponent)
}

func DecodeMockComponent(id string, d func(interface{}) error) (component.Component, error) {
	c := NewComponent(id, nil, nil)

	if err := d(&c); err != nil {
		return c, fmt.Errorf("Failure to unmarshal product 'mock' : %w", err)
	}

	return c, nil
}
