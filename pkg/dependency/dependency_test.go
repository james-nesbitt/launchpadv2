package dependency_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/dependency/mock"
)

// --- test Functions ---

func Test_MockValidate(t *testing.T) {
	ctx := context.Background()
	hds := []dependency.HasDependencies{
		mock.HasDependencies(
			dependency.Requirements{
				mock.Requirement("test", "Validate requirement", nil),
			},
		),
	}
	pds := []dependency.FullfillsDependencies{
		mock.FullfillsDependencies(mock.Dependency("my-dep", "", nil, nil), nil),
	}

	for _, hd := range hds {
		for _, r := range hd.Requires(ctx) {
			for _, pd := range pds {
				d, err := pd.Provides(ctx, r)

				if err != nil {
					t.Error("Dependency matching failed unexpectedly")
					continue
				}
				if d.Id() != "my-dep" {
					t.Errorf("Dependency mathcing did not return the correct dependency: %+v", d)
				}
			}
		}
	}
}

func Test_MatchRequirements_Success(t *testing.T) {
	ctx := context.Background()
	r := mock.Requirement("mock", "", nil)          // only needed as an argument. fds do all of the work
	testerr := errors.New("we should not get here") // this should not get returned
	fds := []dependency.FullfillsDependencies{
		mock.FullfillsDependencies(nil, dependency.ErrDependencyNotHandled),
		mock.FullfillsDependencies(nil, dependency.ErrDependencyNotHandled),
		mock.FullfillsDependencies(mock.Dependency("this-one", "", nil, nil), nil),         // this handler should get used
		mock.FullfillsDependencies(mock.Dependency("not-this-one", "", nil, nil), testerr), // this handler should not get used
	}

	d, err := dependency.GetRequirementDependency(ctx, r, fds)
	if errors.Is(err, testerr) {
		t.Error("dependency handler was not used")
	}
	if err != nil {
		t.Error("MatchRequirements gave unexpected response")
	}

	if d == nil {
		t.Error("MatchRequirements did not return any dependency")
	}
	if d.Id() != "this-one" {
		t.Errorf("MatchRequirements returned the wrong dependency")
	}
}

func Test_MatchRequirements_ShouldFail(t *testing.T) {
	ctx := context.Background()
	r := mock.Requirement("", "", nil)              // only needed as an argument. fds do all of the work
	testerr := errors.New("we should not get here") // this should not get returned
	fds := []dependency.FullfillsDependencies{
		mock.FullfillsDependencies(nil, dependency.ErrDependencyNotHandled),        // this handler should not get used
		mock.FullfillsDependencies(nil, dependency.ErrDependencyShouldHaveHandled), // this handler should get used, and result in an error
		mock.FullfillsDependencies(nil, nil),                                       // this handler should get used
		mock.FullfillsDependencies(nil, testerr),                                   // this handler should not get used
	}

	_, err := dependency.GetRequirementDependency(ctx, r, fds)
	if !errors.Is(err, dependency.ErrDependencyShouldHaveHandled) {
		t.Error("empty dependency handlers did not give expected ShouldHaveHandled error")
	}
}

func Test_MatchRequirements_Empty(t *testing.T) {
	ctx := context.Background()
	r := mock.Requirement("mock", "", nil)      // only needed as an argument. fds do all of the work
	fds := []dependency.FullfillsDependencies{} // empty list of handlers/fullfillers

	_, err := dependency.GetRequirementDependency(ctx, r, fds)
	if !errors.Is(err, dependency.ErrDependencyNotHandled) {
		t.Error("empty dependency handlers did not give expected NotHandled error")
	}
}
