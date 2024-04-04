package v20

import (
	"github.com/Mirantis/launchpad/pkg/cluster"
)

const (
	ID = "launchpad.mirantis.com/v2.0"
)

func Decode(cl *cluster.Cluster, d func(interface{}) error) error {
	var cs Spec

	if err := d(&cs); err != nil {
		return err
	}

	cl.Hosts = cs.Hosts.hosts()
	cl.Components = cs.Products.products()

	return nil
}

// Spec defines cluster spec.
type Spec struct {
	Hosts    SpecHosts    `yaml:"hosts" validate:"required,min=1"`
	Products SpecProducts `yaml:"products" validate:"required,min=1"`
}
