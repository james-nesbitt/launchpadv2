package v20

import (
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/component"
)

type SpecProducts struct {
	cs component.Components
}

func (scs SpecProducts) products() component.Components {
	return scs.cs
}

func (scs *SpecProducts) UnmarshalYAML(py *yaml.Node) error {
	scs.cs = component.Components{}
	spsc := map[string]yaml.Node{}

	if err := py.Decode(&spsc); err != nil {
		return fmt.Errorf("failed to decode spec products")
	}

	errss := []string{}

	for id, spn := range spsc {
		t := id
		pc := struct {
			Id      string `yaml:"id"`
			Handler string `yaml:"handler"`
		}{}

		if err := spn.Decode(&pc); err != nil {
			slog.Debug("Decode of product gave suspiscious results")
		}

		if pc.Handler != "" {
			t = pc.Handler
		}
		if pc.Id != "" {
			id = pc.Id
		}

		p, err := component.DecodeComponent(t, id, spn.Decode)
		if err != nil {
			errss = append(errss, fmt.Sprintf("%s: %s", id, err.Error()))
			continue
		}

		scs.cs[id] = p
	}

	if len(errss) > 0 {
		return fmt.Errorf("Error decoding Spec products: %s", strings.Join(errss, "\n"))
	}

	return nil
}
