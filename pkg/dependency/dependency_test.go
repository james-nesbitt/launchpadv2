package dependency_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// --- Testing Mock structs ---

type Mock_Requirement struct {
	d string
}

func (mr Mock_Requirement) Describe() string {
	return mr.d
}

type Mock_Dependency struct {
	d string
	e error
}

func (md Mock_Dependency) Describe() string {
	return md.d
}

func (md Mock_Dependency) Met() error {
	return md.e
}

type Mock_HasDependencies struct {
	rs dependency.Requirements
}

func (mhc Mock_HasDependencies) Requires(ctx context.Context) dependency.Requirements {
	return mhc.rs
}

type Mock_ProvidesDependencies struct {
	d dependency.Dependency
	h bool // do we handle this type of request
}

func (mpd Mock_ProvidesDependencies) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, bool) {
	return mpd.d, mpd.h
}

// --- test Functions ---

func Test_DependencyMatch_Valid(t *testing.T) {
	ctx := context.Background()

	hds := []dependency.HasDependencies{
		Mock_HasDependencies{
			rs: dependency.Requirements{
				Mock_Requirement{d: "Validate requirement"},
			},
		},
	}
	pds := []dependency.ProvidesDependencies{
		Mock_ProvidesDependencies{
			d: Mock_Dependency{
				d: "Mock description",
			},
			h: true, // yes we handle this type of request
		},
	}

	ds := dependency.Match(ctx, pds, hds)
	uds := ds.Unmet()

	if len(uds) > 0 {
		t.Errorf("Mock dependency validate resulted in unmet dependencies: %+v", uds)
	}
}
