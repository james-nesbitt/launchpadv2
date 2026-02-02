/*
Package dependency defines ways for components to define their dependencies on other components.
*/
package dependency

import (
	"context"
)

// Dependency response to a Dependency Requirement which can fulfill it.
type Dependency interface {
	// ID uniquely identify the Dependency
	ID() string
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
