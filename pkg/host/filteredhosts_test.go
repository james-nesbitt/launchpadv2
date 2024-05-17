package host_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

func Test_FilteredHostDependencySanity(t *testing.T) {
	ctx := context.Background()

	hs := host.NewHosts(
		host.NewHost("one"),
		host.NewHost("two"),
		host.NewHost("three"),
		host.NewHost("four"),
		host.NewHost("five"),
		host.NewHost("six"),
	)

	hfr := host.NewHostsFilterRequirement(
		"test",
		"testing filtered hosts",
		host.HostsFilter(func(_ context.Context, hs host.Hosts) (host.Hosts, error) {
			return hs, nil
		}),
	)

	var r dependency.Requirement = hfr // prove that it is a valid requirement

	hfrm, ok := r.(host.RequiresFilteredHosts) // mutate it to a Filtered Hosts requirement
	if !ok {
		t.Fatalf("requirements could not get mutated")
	}
	if hfrm == nil {
		t.Fatalf("empty filtered hosts requirements mutate")
	}

	ff := hfrm.HostsFilter(ctx) // retrieve the filter

	df := host.HostsDependencyFilterFactory(hs, ff) // build a host filter factory

	fhs, ferr := df(ctx) // run the factory
	if ferr != nil {
		t.Errorf("error occured filtering: %s", ferr.Error())
	}

	if len(fhs) != 6 {
		t.Error("not enought hosts returned in filter")
	}
}
