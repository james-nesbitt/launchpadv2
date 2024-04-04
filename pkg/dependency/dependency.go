package dependency

import (
	"context"
	"errors"
	"log/slog"
)

var (
	ErrDependencyNotMatched        = errors.New("dependency not met")                         // No handler satisfied this dependency
	ErrDependencyShouldHaveHandled = errors.New("dependency not met but it should have been") // This type of dependency is handled, and should have ben handled, but a failure occured.  Don't try any other handler
	ErrDependencyNotHandled        = errors.New("dependency not handled")                     // This type of dependency is not handled so somebody else should handle it (often requivalent to a nil error)
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
	// List the IDs of Requirement.Requirers() which depend on this dependency
	Requirers() []string
}

// --- Dependency and Requirement helper functions

// MatchRequirements Build a dependency for a requirement
//
//		Find a Handler which can handle the Requirement, and make it produce a Dependency.
//		Associate the Dependency with the Requirement, but return it as well so that it
//		can be collected.
//
//	    NOTE: WE DO NOT TELL THE REQUIREMENT ABOUT THE DEPENDENCY - YOU NEED TO DO THAT
func MatchRequirements(ctx context.Context, r Requirement, fds []FullfillsDependencies) (Dependency, error) {
	for _, fd := range fds {
		if fd == nil {
			continue // stupidity check
		}

		d, err := fd.Provides(ctx, r)
		if err == nil {
			slog.Error("Dependency generator succeeded", slog.Any("handler", fd), slog.Any("requirement", r))
			return d, nil // successful dependency creation
		}
		slog.Error("fulfiller err", slog.Any("error", err), slog.Any("fulfiller", fd))

		if errors.Is(err, ErrDependencyShouldHaveHandled) {
			slog.Debug("Dependency generator failed, stopping", slog.Any("handler", fd), slog.Any("requirement", r))
			return d, err // Handler says that it should have handled it, but couldn't, so return an error
		}

		if !errors.Is(err, ErrDependencyNotHandled) {
			slog.Debug("Skipping, as Dependency generate produced an unknown error", slog.String("error", err.Error()))
		}

		continue // Handler error means we should try others
	}

	return nil, ErrDependencyNotMatched
}
