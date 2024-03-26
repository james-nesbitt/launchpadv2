package dependency

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrDependencyNotMet            = errors.New("dependency not met")
	ErrDependencyShouldHaveHandled = errors.New("dependency not met but it should have been") // @TODO come on, this name sucks
	ErrDependencyNotHandled        = errors.New("dependency not handled")
)

// HasDependencies can determine its dependencies
type HasDependencies interface {
	// Requires a set of Requirements indicating what dependencies are needed
	// - if the set is empty, then no requirements are needed
	Requires(context.Context) Requirements
}

// ProvidesDependencies components that can meet dependency needs
type ProvidesDependencies interface {
	// Provides a dependency for some type of Requirements.
	// - A nil dependency means that the type of requirements isn't handled
	// - if the error .Is(ErrDependencyNotMet) then the dependency could not be met
	// - if the error .Is(ErrDependencyNotHandled) then the Provider doesn't handle such requirementd
	// - if the error .Is(ErrDependencyShouldHaveHandled) then the Provider could not
	//   be met, but it should have. No other provider should try to meet it.
	// - A non-nil dependency with a .Met() error indicates that the dependency is not met,
	//   and other providers will be tried.
	Provides(context.Context, Requirement) Dependency
}

type Requirements []Requirement

type Requirement interface {
	Describe() string
	// Fullfill the Dependency by providing a Dependency object.
	// - The Dependency could by [un] Met(), and have an error
	Fullfill(Dependency) error // @TODO us a chan?
}

type Dependencies []Dependency

// UnMet filter for Dependencies that are not met
func (ds Dependencies) Unmet() []Dependency {
	ud := []Dependency{}

	for _, d := range ds {
		if d == nil {
			continue // stupidity filter
		}
		if err := d.Met(); err != nil {
			ud = append(ud, d)
		}
	}

	return ud
}

type Dependency interface {
	// Describe the Dependency so that we can log audit it
	Describe() string
	// Met or not met. If not met provide an error, else provide nil
	Met() error
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

// ProvideRequirementDependency Build a dependency for a requirement
// - we could just fulfill this, but we leave it to the consumer so that they can mutate it and log/audit
func ProvideRequirementDependency(ctx context.Context, r Requirement, pds []ProvidesDependencies) Dependency {
	for _, pd := range pds {
		if pd == nil {
			continue // stupidity check
		}

		d := pd.Provides(ctx, r)

		if d == nil {
			// Provider does not handle this kind of Requirements (acceptable behaviour)
			continue
		}

		err := d.Met()

		if err == nil {
			slog.Debug("Dependency generator succeeded", slog.Any("handler", pd), slog.Any("requirement", r), slog.Any("dependency", d))
			return d // successful dependency creation
		}
		if errors.Is(err, ErrDependencyShouldHaveHandled) {
			slog.Debug("Dependency generator failed, stopping", slog.Any("handler", pd), slog.Any("requirement", r), slog.Any("dependency", d))
			return d // Handler says that it should have handled it, but couldn't, so return an error
		}
		if errors.Is(err, ErrDependencyNotHandled) {
			slog.Debug("Dependency generator failed, trying others", slog.Any("handler", pd), slog.Any("requirement", r), slog.Any("dependency", d))
			continue // Handler failed to create a dependency, but maybe another can
		}
		if errors.Is(err, ErrDependencyNotMet) {
			// TODO log this?
			continue // Handler decided it shouldn't handle the dependency
		}
	}

	return DependencyNotHandled{
		r: r,
	}
}

// === A dependency struct for an unhandled dependency

type DependencyNotHandled struct {
	r Requirement
}

func NewUnmetDependency(requirement Requirement) Dependency {
	return DependencyNotHandled{
		r: requirement,
	}
}

func (dum DependencyNotHandled) Describe() string {
	return dum.r.Describe()
}

func (dum DependencyNotHandled) Met() error {
	return fmt.Errorf("%w; %s", ErrDependencyNotHandled, dum.Describe())
}
