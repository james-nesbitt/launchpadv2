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
func DecodeComponent(id string, d func(any) error) (component.Component, error) {
	c := ComponentDecode{}
	c.Name = id

	if err := defaults.Set(&c); err != nil {
		return nil, fmt.Errorf("failed to set mock defaults: %w", err)
	}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("failure to decode product '%s' : %w", ComponentType, err)
	}

	return NewComponentFromDecode(c), nil
}

// NewComponentFromDecode build a new Mock component from the decode structs.
func NewComponentFromDecode(cd ComponentDecode) component.Component {
	return &comp{
		name: cd.Name,
	}
}

// ComponentDecode data struct used to allow decoded data to be turned into a Mock Component.
type ComponentDecode struct {
	Name string `yaml:"name" json:"name"`
}
