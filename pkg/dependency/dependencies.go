package dependency

import (
	"context"
	"errors"
)

const (
	EventKeyActivated   = "activate"
	EventKeyDeActivated = "deactivate"
)

var (
	// ErrDependencyWrongType something was given the wrong kind of dependency.
	ErrDependencyWrongType = errors.New("received the wrong type of dependency")
	// ErrDependencyNotMatched No handler satisfied this dependency.
	ErrDependencyNotMatched = errors.New("dependency not met")
	// ErrShouldHaveHandled This type of dependency is handled, and should have ben handled, but a failure occurred.  Don't try any other handler.
	ErrShouldHaveHandled = errors.New("dependency not handled but it should have been")
	// ErrNotHandled type of dependency is not handled so somebody else should handle it (often requivalent to a nil error).
	ErrNotHandled = errors.New("dependency not handled by any provider")
)

// FullfillsDependencies components that can meet dependency needs.
type ProvidesDependencies interface {
	// ProvidesDependency a dependency for some type of Requirements.
	// - if the error .Is(ErrDependencyNotHandled) then  Provider doesn't handle such requirement
	// - if the error .Is(ErrDependencyShouldHaveHandled) then the Provider could not
	//   be met, but it should have. No other provider should try to meet it.
	ProvidesDependencies(context.Context, Requirement) (Dependency, error)
}

// Dependencies set of Dependency items.
type Dependencies []Dependency

// NewDependencies a new Dependencies, optionally pre-populated.
func NewDependencies(dsa ...Dependency) Dependencies {
	ds := Dependencies{}
	if len(dsa) > 0 {
		ds = append(ds, dsa...)
	}
	return ds
}

func (ds *Dependencies) Add(nds ...Dependency) {
	*ds = append(*ds, nds...)
}

func (ds *Dependencies) Merge(nds Dependencies) {
	*ds = append(*ds, nds...)
}

func (ds Dependencies) Get(id string) Dependency {
	for _, d := range ds {
		if d.Id() == id {
			return d
		}
	}
	return nil
}

func (ds Dependencies) Ids() []string {
	ids := []string{}
	for _, d := range ds {
		ids = append(ids, d.Id())
	}
	return ids
}

// UnMet dependency requirements set.
func (ds Dependencies) Invalid(ctx context.Context) Dependencies {
	uds := Dependencies{}

	for _, d := range ds {
		if err := d.Validate(ctx); err != nil {
			uds = append(uds, d)
		}
	}

	return uds
}

// Intersect to Dependencies, where two dependencies match if their Id() values match.
func (ds Dependencies) Intersect(ds2 Dependencies) Dependencies {
	icd := Dependencies{} // command dependencies that we own
	for _, d2 := range ds2 {
		for _, d1 := range ds {
			if d1.Id() == d2.Id() {
				icd = append(icd, d1)
			}
		}
	}
	return icd
}
