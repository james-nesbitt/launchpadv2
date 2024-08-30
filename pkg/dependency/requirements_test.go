package dependency_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

// --- Tests ---

func Test_NewRequirements(t *testing.T) {
	r0 := mock.StaticRequirement("zero", "zero", nil)
	r1 := mock.StaticRequirement("one", "one", nil)
	r2 := mock.StaticRequirement("two", "two", nil)
	r3 := mock.StaticRequirement("three", "three", nil)
	r4 := mock.StaticRequirement("four", "four", nil)

	rs := dependency.NewRequirements(r0, r1, r2, r3, r4)

	assert.Len(t, rs, 5, "New requirements had wrong length")

	assert.Same(t, r0, rs[0], "wrong requirement in position")
	assert.Same(t, r1, rs[1], "wrong requirement in position")
	assert.Same(t, r2, rs[2], "wrong requirement in position")
	assert.Same(t, r3, rs[3], "wrong requirement in position")
	assert.Same(t, r4, rs[4], "wrong requirement in position")
}

func Test_CollectRequirements(t *testing.T) {
	ctx := context.Background()
	rs := dependency.Requirements{
		mock.StaticRequirement("zero", "zero", nil),
		mock.StaticRequirement("one", "one", nil),
		mock.StaticRequirement("two", "two", nil),
		mock.StaticRequirement("three", "three", nil),
		mock.StaticRequirement("four", "four", nil),
	}

	hds := []dependency.RequiresDependencies{
		mock.RequiresDependencies(rs[0:2]),                      // the first two
		mock.RequiresDependencies(dependency.NewRequirements()), // empty
		mock.RequiresDependencies(rs[2:3]),                      // the third
		mock.RequiresDependencies(rs[3:]),                       // the rest
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
