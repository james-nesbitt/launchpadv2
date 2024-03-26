package component

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/phase"
)

// Component a cluster component which can provide phase actions
type Component interface {
	// Name of the Component, which is used for logging and auditing
	Name() string

	// Actions to run for a particular string phase
	Actions(context.Context, string) (phase.Actions, error)
}
