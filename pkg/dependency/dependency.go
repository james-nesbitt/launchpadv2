package dependency

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrDependencyWrongType         = errors.New("recieved the wrong type of dependency")          // something wads given the wrong kind of dependency
	ErrDependencyNotMatched        = errors.New("dependency not met")                             // No handler satisfied this dependency
	ErrDependencyShouldHaveHandled = errors.New("dependency not handled but it should have been") // This type of dependency is handled, and should have ben handled, but a failure occured.  Don't try any other handler
	ErrDependencyNotHandled        = errors.New("dependency not handled by any provider")         // This type of dependency is not handled so somebody else should handle it (often requivalent to a nil error)
)

// FullfillsDependencies components that can meet dependency needs
type FullfillsDependencies interface {
	// Provides a dependency for some type of Requirements.
	// - if the error .Is(ErrDependencyNotHandled) then  Provider doesn't handle such requirement
	// - if the error .Is(ErrDependencyShouldHaveHandled) then the Provider could not
	//   be met, but it should have. No other provider should try to meet it.
	Provides(context.Context, Requirement) (Dependency, error)
}

// Dependencies set of Dependency items
type Dependencies []Dependency

func (ds Dependencies) Ids() []string {
	ids := []string{}
	for _, d := range ds {
		ids = append(ids, d.Id())
	}
	return ids
}

// UnMet dependency requirements set
func (ds Dependencies) Invalid(ctx context.Context) Dependencies {
	uds := Dependencies{}

	for _, d := range ds {
		if err := d.Validate(ctx); err != nil {
			uds = append(uds, d)
		}
	}

	return uds
}

// Dependency response to a Dependency Requirement which can fullfill it
type Dependency interface {
	// Id uniquely identify the Dependency
	Id() string
	// Describe the dependency for logging and auditing
	Describe() string // Validate the the dependency thinks it has what it needs to fullfill it
	// Validate the the dependency thinks it has what it needs to fullfill it
	//   Validation does not meet the dependency, it only confirms that the dependency can be met
	//   when it is needed.  Each requirement will have its own interface for individual types of
	//   requirements.
	Validate(context.Context) error
	// Met is the dependency fullfilled, or is it still blocked/waiting for fulfillment
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
func GetRequirementDependency(ctx context.Context, r Requirement, fds []FullfillsDependencies) (Dependency, error) {
	if len(fds) == 0 {
		return nil, fmt.Errorf("%w; no dependency handlers providers for requirement %s", ErrDependencyNotHandled, r.Id())
	}

	for _, fd := range fds {
		if fd == nil {
			continue // stupidity check
		}

		d, err := fd.Provides(ctx, r)

		// it is allowed to return no error, and no dependency,
		// which just means that the req wasn't handled
		if d == nil && err == nil {
			err = ErrDependencyNotHandled
		}

		if err == nil {
			return d, nil // successful dependency creation
		}

		if errors.Is(err, ErrDependencyShouldHaveHandled) {
			return d, err // Handler says that it should have handled it, but couldn't, so return an error
		}

		if !errors.Is(err, ErrDependencyNotHandled) {
			slog.WarnContext(ctx, fmt.Sprintf("dependency: Unknown Dependency generation error: (%s) %s", r.Id(), err.Error()), slog.String("error", err.Error()))
		}
	}

	return nil, ErrDependencyNotHandled
}
