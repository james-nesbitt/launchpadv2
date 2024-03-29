package dependency_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// --- Testing Mock structs ---

type Mock_HasDependencies struct {
	rs dependency.Requirements
}

func (mhc Mock_HasDependencies) Requires(ctx context.Context) dependency.Requirements {
	return mhc.rs
}

type Mock_ProvidesDependencies struct {
	err error
}

func (mpd Mock_ProvidesDependencies) Provides(ctx context.Context, r dependency.Requirement) error {
	return mpd.err
}

type Mock_Requirement struct {
	description string
	err         error
}

func (mr Mock_Requirement) Describe() string {
	return mr.description
}

func (mr Mock_Requirement) Validate(_ context.Context) error {
	return mr.err
}

// --- test Functions ---

// Not much of a test, the only coverage is the Requirements.UnMet(), but
// it provides a general functional test of the dependency matching approach
// using mocks.
func Test_DependencyMatch_Valid(t *testing.T) {
	ctx := context.Background()

	hds := []dependency.HasDependencies{
		Mock_HasDependencies{
			rs: dependency.Requirements{
				Mock_Requirement{
					description: "Validate requirement",
					err:         nil,
				},
			},
		},
	}
	pds := []dependency.FullfillsDependencies{
		Mock_ProvidesDependencies{
			err: nil,
		},
	}

	rs := dependency.Requirements{}

	for _, hd := range hds {
		for _, r := range hd.Requires(ctx) {
			rs = append(rs, r)

			for _, pd := range pds {
				err := pd.Provides(ctx, r)

				if err != nil {
					t.Error("Dependency matching failed unexpectedly")
				}
			}
		}
	}

	urs := rs.UnMet(ctx)

	if len(urs) > 0 {
		t.Errorf("Mock dependency validate resulted in unmet dependencies: %+v", urs)
	}
}

func Test_CollectRequirements(t *testing.T) {
	ctx := context.Background()
	rs := dependency.Requirements{
		Mock_Requirement{
			description: "zero",
		},
		Mock_Requirement{
			description: "one",
		},
		Mock_Requirement{
			description: "two",
		},
		Mock_Requirement{
			description: "three",
		},
		Mock_Requirement{
			description: "four",
		},
	}

	hds := []dependency.HasDependencies{
		Mock_HasDependencies{
			rs: rs[0:2], // the first two
		},
		Mock_HasDependencies{ // empty
			rs: dependency.Requirements{},
		},
		Mock_HasDependencies{
			rs: rs[2:3], // the third
		},
		Mock_HasDependencies{
			rs: rs[3:], // the rest
		},
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

func Test_FullfillRequirements_Success(t *testing.T) {
	ctx := context.Background()
	r := Mock_Requirement{}                         // only needed as an argument. fds do all of the work
	testerr := errors.New("we should not get here") // this should not get returned
	fds := []dependency.FullfillsDependencies{
		Mock_ProvidesDependencies{
			err: dependency.ErrDependencyNotHandled,
		},
		Mock_ProvidesDependencies{
			err: dependency.ErrDependencyNotHandled,
		},
		Mock_ProvidesDependencies{ // this handler should get used
			err: nil,
		},
		Mock_ProvidesDependencies{ // this handler should not get used
			err: testerr,
		},
	}

	err := dependency.FullfillRequirements(ctx, r, fds)
	if errors.Is(err, testerr) {
		t.Error("dependency handler was not used")
	}
	if err != nil {
		t.Error("FullfillRequirments gave unexpected response")
	}
}

func Test_FullfillRequirements_ShouldFail(t *testing.T) {
	ctx := context.Background()
	r := Mock_Requirement{}                         // only needed as an argument. fds do all of the work
	testerr := errors.New("we should not get here") // this should not get returned
	fds := []dependency.FullfillsDependencies{
		Mock_ProvidesDependencies{ // this handler should not get used
			err: dependency.ErrDependencyNotHandled,
		},
		Mock_ProvidesDependencies{ // this handler should get used, and result in an error
			err: dependency.ErrDependencyShouldHaveHandled,
		},
		Mock_ProvidesDependencies{ // this handler should not get used
			err: nil,
		},
		Mock_ProvidesDependencies{ // this handler should not get used
			err: testerr,
		},
	}

	err := dependency.FullfillRequirements(ctx, r, fds)
	if !errors.Is(err, dependency.ErrDependencyShouldHaveHandled) {
		t.Error("empty dependency handlers did not give expected NotMet error")
	}
}

func Test_FullfillRequirements_Empty(t *testing.T) {
	ctx := context.Background()
	r := Mock_Requirement{}                     // only needed as an argument. fds do all of the work
	fds := []dependency.FullfillsDependencies{} // empty list of handlers/fullfillers

	err := dependency.FullfillRequirements(ctx, r, fds)
	if !errors.Is(err, dependency.ErrDependencyNotMet) {
		t.Error("empty dependency handlers did not give expected NotMet error")
	}
}
