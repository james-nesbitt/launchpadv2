package host_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	mockhost "github.com/Mirantis/launchpad/pkg/host/mock"

	"github.com/Mirantis/launchpad/pkg/host"
)

func Test_HostComponent(t *testing.T) {
	id := "cluster1"
	hs := host.Hosts{
		"one":   mockhost.NewHost("one", []string{"manager"}, nil),
		"two":   mockhost.NewHost("two", []string{"worker"}, nil),
		"three": mockhost.NewHost("three", []string{"worker"}, nil),
	}

	var hc component.Component = host.NewHostsComponent(id, hs)

	if !strings.Contains(hc.Name(), id) {
		t.Errorf("Hosts component gave wrong name: %s", hc.Name())
	}
}

func Test_HostComponentDependencies(t *testing.T) {
	ctx := context.Background()
	id := "cluster1"
	hs := host.Hosts{
		"zero":  mockhost.NewHost("zero", []string{"ignoreme"}, nil),
		"one":   mockhost.NewHost("one", []string{"manager"}, nil),
		"two":   mockhost.NewHost("two", []string{"worker"}, nil),
		"three": mockhost.NewHost("three", []string{"worker"}, nil),
	}

	hc := host.NewHostsComponent(id, hs)

	var hcpd dependency.ProvidesDependencies = hc

	r := host.NewHostsRolesRequirement("dummy", "dummy hosts requirement", host.HostsRolesFilter{Roles: []string{"manager", "worker"}, Min: 1, Max: 0})

	d, err := hcpd.ProvidesDependencies(context.Background(), r)
	if err != nil {
		t.Errorf("host component Provides returned an error on empty requirement: %s : %+v", err.Error(), d)
	}

	phd, ok := d.(host.HostsDependency)
	if !ok {
		t.Fatalf("hosts roles requirement is wrong type: %+v", d)
	}

	hs2 := phd.ProduceHosts(ctx)
	if len(hs2) != 3 {
		t.Errorf("wrong number of hosts returned: %+v", hs2)
	}

	if hs2["one"].Id() != hs["one"].Id() {
		t.Errorf("wrong first filtered host: %+v", hs2)
	}
}
