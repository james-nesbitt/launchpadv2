package docker

type Version struct {
	Version string
	Plugins []string
}

// NewDocker Implementation.
func NewDocker() *Docker {
	return &Docker{}
}

// Docker host level Docker implementation.
type Docker struct {
	state State
}

// State for the docker implementation on hosts.
type State struct {
}
