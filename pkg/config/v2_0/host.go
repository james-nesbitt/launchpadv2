package v20

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	ctx, c := context.WithTimeout(context.Background(), time.Second*30)
	defer c()

	shs.hs = host.NewHosts()

	// decode into a map of hosts, where each host is a map of host plugins
	shsc := map[string]map[string]yaml.Node{}
	if err := py.Decode(&shsc); err != nil {
		return fmt.Errorf("failed to decode spec hosts: %s", err.Error())
	}

	if len(shsc) == 0 {
		return fmt.Errorf("no hosts were found when decoding the project")
	}

	errs := []error{}
	for id, shc := range shsc {
		ds := map[string]func(interface{}) error{}
		for k, d := range shc { // collect all of the plugin node .Decode functions
			ds[k] = d.Decode
		}

		h, err := host.DecodeHost(ctx, id, ds)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", id, err.Error()))
			continue
		}

		shs.hs.Add(h)
	}

	if len(errs) > 0 {
		return fmt.Errorf("2.0 host: Error decoding Spec hosts: %s", errors.Join(errs...).Error())
	}

	return nil
}
