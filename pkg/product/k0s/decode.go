package k0s

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeComponent)
}

// DecodeComponent decode a new component from an unmarshall decoder.
func DecodeComponent(id string, d func(interface{}) error) (component.Component, error) {
	var c Config

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponent(id, c), nil
}
