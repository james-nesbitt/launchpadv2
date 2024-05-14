package v20

import (
	"log/slog"
	"strings"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/config"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	ID = "launchpad.mirantis.com/v2.0"
)

func init() {
	config.RegisterSpecDecoder(ID, Decode)
}

// Decode a project from the spec.
func Decode(cl *project.Project, d func(interface{}) error) error {
	var cs Spec

	if err := d(&cs); err != nil {
		return err
	}

	if len(cs.Hosts.hs) == 0 {
		slog.Warn("no hosts were detected")
	}
	if cs.Products.cs == nil {
		slog.Warn("no components were detected")
	}

	if cl.Components == nil {
		cl.Components = component.Components{}
	}

	cl.Components[host.ComponentType] = host.NewHostsComponent(host.ComponentType, cs.Hosts.hosts())

	for id, cp := range cs.Products.products() {
		cl.Components[strings.ToUpper(id)] = cp
	}

	return nil
}
