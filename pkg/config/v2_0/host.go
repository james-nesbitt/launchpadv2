package v20

import (
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/host"
)

type SpecHosts struct {
	hs host.Hosts
}

func (shs SpecHosts) hosts() host.Hosts {
	return shs.hs
}

func (shs *SpecHosts) UnmarshalYAML(py *yaml.Node) error {
	shs.hs = host.Hosts{}
	shsc := map[string]yaml.Node{}

	if err := py.Decode(&shsc); err != nil {
		return fmt.Errorf("failed to decode spec hosts: %s", err.Error())
	}

	errss := []string{}

	for id, shn := range shsc {
		t := id
		pc := struct {
			Id      string `yaml:"id"`
			Handler string `yaml:"handler"`
		}{}

		if err := shn.Decode(&pc); err != nil {
			slog.Debug("Decode of host gave suspiscious results")
		}

		if pc.Handler != "" {
			t = pc.Handler
		}
		if pc.Id != "" {
			id = pc.Id
		}

		p, err := host.DecodeHost(t, id, shn.Decode)
		if err != nil {
			errss = append(errss, fmt.Sprintf("%s: %s", id, err.Error()))
			continue
		}

		shs.hs[id] = p
	}

	if len(errss) > 0 {
		return fmt.Errorf("Error decoding Spec hosts: %s", strings.Join(errss, "\n"))
	}

	return nil
}
