package dependency

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
type FullfillsDependencies interface {
	// ProvidesDependency a dependency for some type of Requirements.
	// - if the error .Is(ErrDependencyNotHandled) then  Provider doesn't handle such requirement
	// - if the error .Is(ErrDependencyShouldHaveHandled) then the Provider could not
	//   be met, but it should have. No other provider should try to meet it.
	Provides(context.Context, Requirement) (Dependency, error)
}

// Dependencies set of Dependency items.
type Dependencies []Dependency

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

// Dependency response to a Dependency Requirement which can fulfill it.
type Dependency interface {
	// Id uniquely identify the Dependency
	Id() string
	// Describe the dependency for logging and auditing
	Describe() string // Validate the the dependency thinks it has what it needs to fulfill it
	// Validate the the dependency thinks it has what it needs to fulfill it
	//   Validation does not meet the dependency, it only confirms that the dependency can be met
	//   when it is needed.  Each requirement will have its own interface for individual types of
	//   requirements.
	Validate(context.Context) error
	// Met is the dependency fulfilled, or is it still blocked/waiting for fulfillment
	Met(context.Context) error
}

// --- Dependency and Requirement helper functions

// GetRequirementDependency Build a dependency for a requirement
//
//		Find a Handler which can handle the Requirement, and make it produce a Dependency.
//		Associate the Dependency with the Requirement, but return it as well so that it
//		can be collected.
//
//	    NOTE: WE DO NOT TELL THE REQUIREMENT ABOUT THE DEPENDENCY - YOU NEED TO DO THAT
func GetRequirementDependency(ctx context.Context, r Requirement, fds map[string]FullfillsDependencies) (Dependency, error) {
	if len(fds) == 0 {
		return nil, fmt.Errorf("%w; no dependency handlers providers for requirement %s", ErrNotHandled, r.Id())
	}

	for name, fd := range fds {
		if fd == nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s handler empty", name), slog.Any("requirement", r), slog.Any("handler", fd))
			continue // stupidity check
		}

		d, err := fd.Provides(ctx, r)

		// it is allowed to return no error, and no dependency,
		// which just means that the req wasn't handled
		if d == nil && err == nil {
			err = ErrNotHandled
		}

		if err == nil {
			slog.DebugContext(ctx, fmt.Sprintf("dependency: '%s' handled by '%s'", r.Id(), name), slog.Any("requirement", r), slog.Any("handler", fd))
			return d, nil // successful dependency creation
		}

		if errors.Is(err, ErrShouldHaveHandled) {
			slog.WarnContext(ctx, fmt.Sprintf("dependency: Dependency generation failure: (%s) %s", r.Id(), err.Error()), slog.String("error", err.Error()))
			return d, err // Handler says that it should have handled it, but couldn't, so return an error
		}

		if !errors.Is(err, ErrNotHandled) {
			slog.WarnContext(ctx, fmt.Sprintf("dependency: Unknown Dependency generation error: (%s) %s", r.Id(), err.Error()), slog.String("error", err.Error()))
		}
	}

	return nil, ErrNotHandled
}
