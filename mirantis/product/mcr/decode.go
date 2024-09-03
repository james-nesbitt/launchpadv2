package mcr

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/creasty/defaults"
)

func init() {
	component.RegisterDecoder("mcr", DecodeMCRComponent)

}

func DecodeMCRComponent(id string, d func(interface{}) error) (component.Component, error) {
	c := Config{}

	defaults.Set(&c)

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MCR' : %w", err)
	}

	return NewComponent(id, c), nil
}
