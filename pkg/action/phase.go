/*
Package action collections of functionality that can be executed as a single run.
*/
package action

import (
	"context"
)

// Phase an executable action in a command
//
//	NOTE that Phases get sorted based on events, but you can use the events
//	  interfaces for that.
type Phase interface {
	// Uniquely identify the phase
	ID() string
	// Execute the phase
	Run(context.Context) error
}
