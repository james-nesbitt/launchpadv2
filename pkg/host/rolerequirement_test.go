package host_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

// --- Tests ---

func Test_ClusterHostsRoleRequirementStandard_Sanity(t *testing.T) {
	ctx := context.Background()
	hs := host.Hosts{
		TestingHost{
			roles: []string{"one"},
		},
	}

	chrs := host.ClusterHostsRoleRequirementStandard{
		Description:  "test",
		Role:         "one",
		MinimumCount: 1,
	}

	// simulate what a dependency handler will do with the  requirement
	// 1. they only accept Requirement interface
	// 2. they will check to see if it can be converted to a RolesRequirement interface
	// 3. a master host list will be assigned
	// 4. matching hosts should be retrievable from the struct (by its provider)

	var r dependency.Requirement = &chrs
	chrr, ok := r.(host.ClusterHostsRequirement)
	if !ok {
		t.Errorf("Could not convert ClusterHostsRoleRequirementStandard to a ClusterHostsRoleRequirement: %+v", r)
	}
	chrr.AssignSourceHosts(hs)

	if err := r.Validate(ctx); err != nil {
		t.Errorf("Could not validate: %s", err.Error())
	}
}

func Test_ClusterHostsRoleRequirementStandard_MinMax(t *testing.T) {
	ctx := context.Background()
	hs := host.Hosts{
		TestingHost{
			roles: []string{"one"},
		},
		TestingHost{
			roles: []string{"two"},
		},
		TestingHost{
			roles: []string{"two", "four"},
		},
		TestingHost{
			roles: []string{"two"},
		},
		TestingHost{
			roles: []string{"two"},
		},
		TestingHost{
			roles: []string{"four"},
		},
		TestingHost{
			roles: []string{"five"},
		},
		TestingHost{
			roles: []string{"six"},
		},
		TestingHost{
			roles: []string{"one"},
		},
	}

	chrs1 := host.ClusterHostsRoleRequirementStandard{
		Description:  "gve me the first three 'two' hosts",
		Role:         "two",
		MinimumCount: 2,
		MaximumCount: 3,
	}
	var r1 dependency.Requirement = &chrs1
	chrr1, ok := r1.(host.ClusterHostsRequirement)
	if !ok {
		t.Errorf("Could not convert ClusterHostsRoleRequirementStandard to a ClusterHostsRoleRequirement: %+v", r1)
	}
	chrr1.AssignSourceHosts(hs)
	if err := r1.Validate(ctx); err != nil {
		t.Errorf("Could not validate1: %s", err.Error())
	}

	chrs2 := host.ClusterHostsRoleRequirementStandard{
		Description:  "there are not enough fives",
		Role:         "five",
		MinimumCount: 2,
	}
	var r2 dependency.Requirement = &chrs2
	chrr2, ok := r2.(host.ClusterHostsRequirement)
	if !ok {
		t.Errorf("Could not convert ClusterHostsRoleRequirementStandard to a ClusterHostsRoleRequirement: %+v", r2)
	}
	chrr2.AssignSourceHosts(hs)
	if err := r2.Validate(ctx); err == nil {
		t.Errorf("Validation2 passed when it was expected to fail: %+s", err.Error())
	}

	chrs3 := host.ClusterHostsRoleRequirementStandard{
		Description:  "give me two 'four's",
		Role:         "four",
		MinimumCount: 2,
		MaximumCount: 2,
	}
	var r3 dependency.Requirement = &chrs3
	chrr3, ok := r3.(host.ClusterHostsRequirement)
	if !ok {
		t.Errorf("Could not convert ClusterHostsRoleRequirementStandard to a ClusterHostsRoleRequirement: %+v", r3)
	}
	chrr3.AssignSourceHosts(hs)
	if err := r3.Validate(ctx); err != nil {
		t.Errorf("Could not validate3: %s", err.Error())
	}

	if mhs3, err := chrs3.RetrieveMatchingHosts(); err != nil {
		t.Errorf("Retrieving hosts produced an error: %s", err.Error())
	} else if len(mhs3) != 2 {
		t.Errorf("The wrong number of hosts was returned: %+v", mhs3)
	}
}
