package v20

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/host"
)

type SpecHosts map[string]*SpecHost

func (sc SpecHosts) hosts() host.Hosts {
	hs := host.Hosts{}

	for id, h := range sc {
		hs[id] = h.host
	}

	return hs
}

type SpecHost struct {
	host host.Host
}

func (sh *SpecHost) UnmarshalYAML(py *yaml.Node) error {
	t := "rig" // hosts will be rig if not otherwise stated
	hc := struct {
		Handler string `yaml:"handler"`
		Role    string `yaml:"role"`
	}{}

	if err := py.Decode(&hc); err != nil {
		slog.Error("host decode gave suspicious contents")
	}

	if hc.Handler != "" {
		t = hc.Handler
	}

	h, err := host.DecodeHost(t, py.Decode)
	if err != nil {
		return fmt.Errorf("host decode failed: %s", err.Error())
	}

	sh.host = h
	return nil
}
