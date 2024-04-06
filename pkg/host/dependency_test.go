package host_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/mock"
)

func Test_DependencyHostRolesSanity(t *testing.T) {
	ctx := context.Background()
	hs := host.Hosts{
		"mock": mock.NewHost("mock", []string{"dummy"}, nil),
	}
	hrr := host.NewHostsRolesRequirement("test", "my test hostroles req", host.HostsRolesFilter{Roles: []string{"dummy"}, Min: 1, Max: 1})
	hd := host.NewHostsDependency("test", "my test hostroles dep", func(ctx context.Context) (*host.Hosts, error) {
		return &hs, nil
	})

	var r dependency.Requirement = hrr
	var d dependency.Dependency = hd

	if _, ok := d.(host.HostsDependency); !ok {
		t.Errorf("our requirement type doesn't match our requirement interface")
	}
	if _, ok := d.(host.HostsDependency); !ok {
		t.Errorf("our dependency type doesn't match our dependency interface")
	}

	if dm := r.Matched(ctx); dm != nil {
		t.Errorf("our requirements says it is matched before we matched: %+v", dm)
	}

	if err := r.Match(d); err != nil {
		t.Errorf("requirement match returned an error: %s", err.Error())
	}

	d2 := r.Matched(ctx)
	hd2, ok := d2.(host.HostsDependency)

	if !ok {
		t.Errorf("couldn't convert requirement dependency to our interface")
	}

	hs2p := hd2.ProduceHosts(ctx)
	hs2 := *hs2p

	if len(hs2) != 1 {
		t.Errorf("host dpeendnecy didn't return the right number of hosts")
	}
}
