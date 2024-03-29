package dependency

import (
	"context"
	"errors"
	"log/slog"
)

var (
	ErrDependencyNotMet            = errors.New("dependency not met")                         // No handler satisfied this dependency
	ErrDependencyShouldHaveHandled = errors.New("dependency not met but it should have been") // This type of dependency is handled, and should have ben handled, but a failure occured.  Don't try any other handler
	ErrDependencyNotHandled        = errors.New("dependency not handled")                     // This type of dependency is not handled so somebody else should handle it (often requivalent to a nil error)
)

// HasDependencies can determine its dependencies
type HasDependencies interface {
	// Requires a set of Requirements indicating what dependency Requirements are needed
	// - if the set is empty, then no requirements are needed
	Requires(context.Context) Requirements
}

// FullfillsDependencies components that can meet dependency needs
type FullfillsDependencies interface {
	// Provides a dependency for some type of Requirements.
	// - A nil return means that the type of requirements isn't handled
	// - if the error .Is(ErrDependencyNotMet) then the dependency could not be met
	// - if the error .Is(ErrDependencyNotHandled) then the Provider doesn't handle such requirementd
	// - if the error .Is(ErrDependencyShouldHaveHandled) then the Provider could not
	//   be met, but it should have. No other provider should try to meet it.
	// - A non-nil dependency with a .Met() error indicates that the dependency is not met,
	//   and other providers will be tried.
	Provides(context.Context, Requirement) error
}

// --- Dependency and Requirement helper functions

// CollectRequirements from a set of HasDependencies
func CollectRequirements(ctx context.Context, hds []HasDependencies) Requirements {
	rs := Requirements{}
	for _, hd := range hds {
		if hd == nil {
			continue // stupidity check
		}

		for _, r := range hd.Requires(ctx) {
			slog.Debug("Dependency requirement received", slog.Any("handler", hd), slog.Any("requirement", r))
			rs = append(rs, r)
		}
	}
	return rs
}

// FullfillRequirements Build a dependency for a requirement
// - we could just fulfill this, but we leave it to the consumer so that they can mutate it and log/audit
func FullfillRequirements(ctx context.Context, r Requirement, fds []FullfillsDependencies) error {
	for _, fd := range fds {
		if fd == nil {
			continue // stupidity check
		}

		err := fd.Provides(ctx, r)

		if err == nil {
			slog.Debug("Dependency generator succeeded", slog.Any("handler", fd), slog.Any("requirement", r))
			return nil // successful dependency creation
		}
		if errors.Is(err, ErrDependencyShouldHaveHandled) {
			slog.Debug("Dependency generator failed, stopping", slog.Any("handler", fd), slog.Any("requirement", r))
			return err // Handler says that it should have handled it, but couldn't, so return an error
		}
		if errors.Is(err, ErrDependencyNotHandled) {
			slog.Debug("Dependency generator failed, trying others", slog.Any("handler", fd), slog.Any("requirement", r))
			continue // Handler failed to create a dependency, but maybe another can
		}
		if errors.Is(err, ErrDependencyNotMet) {
			// TODO log this?
			continue // Handler decided it shouldn't handle the dependency
		}
	}

	return ErrDependencyNotMet
}
