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
	pds := map[string]dependency.FullfillsDependencies{
		"my-dep": mock.FullfillsDependencies(mock.Dependency("my-dep", "", nil, nil), nil),
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
	fds := map[string]dependency.FullfillsDependencies{
		"skip-one":     mock.FullfillsDependencies(nil, dependency.ErrNotHandled),
		"skip-two":     mock.FullfillsDependencies(nil, dependency.ErrNotHandled),
		"use-this-one": mock.FullfillsDependencies(mock.Dependency("this-one", "", nil, nil), nil),         // this handler should get used
		"skip-three":   mock.FullfillsDependencies(mock.Dependency("not-this-one", "", nil, nil), testerr), // this handler should not get used
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
	fds := map[string]dependency.FullfillsDependencies{
		"skip-one":             mock.FullfillsDependencies(nil, dependency.ErrNotHandled),        // this handler should not get used
		"should-fail-handling": mock.FullfillsDependencies(nil, dependency.ErrShouldHaveHandled), // this handler should get used, and result in an error
		"skip-two":             mock.FullfillsDependencies(nil, nil),                             // this handler should get used
		"skip-three":           mock.FullfillsDependencies(nil, testerr),                         // this handler should not get used
	}

	_, err := dependency.GetRequirementDependency(ctx, r, fds)
	if !errors.Is(err, dependency.ErrShouldHaveHandled) {
		t.Error("empty dependency handlers did not give expected ShouldHaveHandled error")
	}
}

func Test_MatchRequirements_Empty(t *testing.T) {
	ctx := context.Background()
	r := mock.Requirement("mock", "", nil)               // only needed as an argument. fds do all of the work
	fds := map[string]dependency.FullfillsDependencies{} // empty list of handlers/fullfillers

	_, err := dependency.GetRequirementDependency(ctx, r, fds)
	if !errors.Is(err, dependency.ErrNotHandled) {
		t.Error("empty dependency handlers did not give expected NotHandled error")
	}
}

func Test_DependenciesSanity(t *testing.T) {
	d1 := mock.Dependency("one", "", nil, nil)
	d2 := mock.Dependency("two", "", nil, nil)
	d3 := mock.Dependency("three", "", nil, nil)
	d4 := mock.Dependency("four", "", nil, nil)
	d5 := mock.Dependency("five", "", nil, nil)
	d6 := mock.Dependency("six", "", nil, nil)
	d7 := mock.Dependency("seven", "", nil, nil)

	dsA := dependency.Dependencies{d1, d2, d3, d4, d6}
	dsB := dependency.Dependencies{d3, d4, d5, d5}
	dsC := dependency.Dependencies{}

	if len(dsC.Ids()) != 0 {
		t.Errorf("empty dependencies gave non-empty Id(): %+v", dsC)
	}

	if d := dsB.Get(d4.Id()); d == nil {
		t.Errorf("dependencies Get for known id returned null: %+v", dsB)
	} else if d.Id() != d4.Id() {
		t.Errorf("dependencies get for known id has the wrong id: %+v", dsA)
	}

	dsA.Add(d7)

	if d := dsA.Get(d7.Id()); d == nil {
		t.Errorf("dependencies get for added id returned nil: %+v", dsA)
	} else if d.Id() != d7.Id() {
		t.Errorf("dependencies get for added id has the wrong id: %+v", dsA)
	}
}

func Test_DependencyIntersect(t *testing.T) {
	d1 := mock.Dependency("one", "", nil, nil)
	d2 := mock.Dependency("two", "", nil, nil)
	d3 := mock.Dependency("three", "", nil, nil)
	d4 := mock.Dependency("four", "", nil, nil)
	d5 := mock.Dependency("five", "", nil, nil)
	d6 := mock.Dependency("six", "", nil, nil)

	dsA := dependency.Dependencies{d1, d2, d3, d4, d6}
	dsB := dependency.Dependencies{d3, d4, d5, d5}
	dsC := dependency.Dependencies{}

	dsiOne := dsA.Intersect(dsB)
	if len(dsiOne) != 2 {
		t.Errorf("Incorrect dependencies intersect: %+v", dsiOne)
	}

	dsiTwo := dsA.Intersect(dsC)
	if len(dsiTwo) != 0 {
		t.Errorf("Incorrect dependencenies for empty participant: %+v", dsiTwo)
	}

}
