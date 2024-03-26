package docker

import "context"

// ProvidesImplementation indicates that it can create a Docker ImplementationDocker
type ProvidesImplementation interface {
	ImplementDocker(context.Context) (Docker, error)
}

// NewDocker Implementation
func NewDocker() Docker {
	return Docker{}
}

// Docker host level Docker implementation
type Docker struct {
	state State
}

// State for the docker implementation on hosts
type State struct {
}
