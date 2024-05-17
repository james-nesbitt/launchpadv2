package dependency_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
)

// --- Tests ---

func Test_CollectRequirements(t *testing.T) {
	ctx := context.Background()
	rs := dependency.Requirements{
		mock.Requirement("zero", "zero", nil),
		mock.Requirement("one", "one", nil),
		mock.Requirement("two", "two", nil),
		mock.Requirement("three", "three", nil),
		mock.Requirement("four", "four", nil),
	}

	hds := []dependency.RequiresDependencies{
		mock.RequiresDependencies(rs[0:2]),                   // the first two
		mock.RequiresDependencies(dependency.Requirements{}), // empty
		mock.RequiresDependencies(rs[2:3]),                   // the third
		mock.RequiresDependencies(rs[3:]),                    // the rest
	}

	crs := dependency.CollectRequirements(ctx, hds)

	if len(crs) != len(rs) {
		t.Errorf("Collect requirements did not return the expected number of requirements: %+v != %+v", crs, rs)
	} else {
		for i, r := range rs {
			if crs[i].Describe() != r.Describe() {
				t.Error("Collect requirments did not return the expected requirements/order")
			}
		}
	}
}
