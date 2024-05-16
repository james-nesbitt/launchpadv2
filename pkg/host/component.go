package host

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "hosts"
)

func NewHostsComponent(id string, hosts Hosts) *HostsComponent {
	return &HostsComponent{
		id:    id,
		hosts: hosts,
		deps:  dependency.Dependencies{},
	}
}

type HostsComponent struct {
	id    string
	hosts Hosts

	deps dependency.Dependencies
}

func (hc HostsComponent) Name() string {
	if hc.id == ComponentType {
		return hc.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, hc.id)
}

func (hc HostsComponent) Debug() interface{} {
	return struct {
		ID    string `json:"ID"`
		Hosts Hosts  `json:"Hosts"`
	}{
		ID:    hc.id,
		Hosts: hc.hosts,
	}
}

func (hc HostsComponent) Validate(_ context.Context) error {
	return nil
}

func (hc HostsComponent) Hosts() Hosts {
	return hc.hosts
}
