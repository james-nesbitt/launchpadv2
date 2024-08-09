package v20

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
)

type SpecHosts struct {
	hc component.Component
}

func (shs SpecHosts) hostComponent() component.Component {
	return shs.hc
}

func (shs *SpecHosts) UnmarshalYAML(py *yaml.Node) error {
	hc, err := component.DecodeComponent(host.ComponentType, host.ComponentType, py.Decode)
	shs.hc = hc

	if err != nil {
		return fmt.Errorf("error decoding hosts")
	}

	return nil
}
