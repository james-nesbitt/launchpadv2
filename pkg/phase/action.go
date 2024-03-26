package phase

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

type Actions []Action

// Action step in a phase
type Action interface {
	// Label for logging and labelling
	Label() string
	// Validate the action
	Validate(context.Context, host.Hosts, Actions) ([]string, error)
	// Run the Action
	Run(context.Context) error
}
