package mock

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/creasty/defaults"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeComponent)
}

// DecodeComponent component decoder used with the component injection system.
func DecodeComponent(id string, d func(interface{}) error) (component.Component, error) {
	c := ComponentDecode{}
	c.Name = id

	defaults.Set(&c)

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponentFromDecode(c), nil
}

// NewComponentFromDecode build a new Mock component from the decode structs.
func NewComponentFromDecode(cd ComponentDecode) component.Component {
	return &comp{
		name: cd.Name,
	}
}

// ComputeDecode data struct used to allow decoded data to be turned into a Mock Component.
type ComponentDecode struct {
	Name string `yaml:"name" json:"name"`
}
