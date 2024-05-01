package v20

import (
	"github.com/Mirantis/launchpad/pkg/cluster"
	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/config"
	"github.com/Mirantis/launchpad/pkg/host"
)

const (
	ID = "launchpad.mirantis.com/v2.0"

	HostsComponentId = "hosts"
)

func init() {
	config.RegisterSpecDecoder(ID, Decode)
}

// Decode a cluster from the spec.
func Decode(cl *cluster.Cluster, d func(interface{}) error) error {
	var cs Spec

	if err := d(&cs); err != nil {
		return err
	}

	if cl.Components == nil {
		cl.Components = component.Components{}
	}

	cl.Components[HostsComponentId] = host.NewHostsComponent(HostsComponentId, cs.Hosts.hosts())

	for id, cp := range cs.Products.products() {
		cl.Components[id] = cp
	}

	return nil
}
