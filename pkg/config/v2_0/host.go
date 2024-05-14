package v20

import (
	"errors"
	"fmt"
	"log/slog"

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

	errs := []error{}

	if len(shsc) == 0 {
		errs = append(errs, fmt.Errorf("no hosts were found when decoding the project"))
	}
	for id, shn := range shsc {
		pc := struct {
			Id      string   `yaml:"id"`
			Handler string   `yaml:"handler"`
			Roles   []string `yaml:"roles"`
		}{}

		if err := shn.Decode(&pc); err != nil {
			slog.Debug("Decode of host gave suspiscious results")
		}

		if pc.Handler == "" {
			errs = append(errs, fmt.Errorf("no 'handler' specified for host '%s'", id))
			continue
		}

		h, err := host.DecodeHost(pc.Handler, id, shn.Decode)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", id, err.Error()))
			continue
		}

		shs.hs[id] = h
	}

	if len(errs) > 0 {
		return fmt.Errorf("2.0 host: Error decoding Spec hosts: %s", errors.Join(errs...).Error())
	}

	return nil
}
