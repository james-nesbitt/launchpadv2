/*
Package component define component plugin functionality for launchpad.
*/
package component

import (
	"context"
)

// Component a project component which can provide phase actions.
type Component interface {
	// Name of the Component, which is used for logging and auditing
	Name() string
	// Debug the component by returning something that can be serialized
	Debug() any
	// Validate the component thinks it has valid config
	Validate(context.Context) error
}
