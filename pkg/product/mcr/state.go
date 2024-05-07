package mcr

import dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"

// state state object for MCR, to be used to compute delta from config.
type state struct {
	leader *dockerhost.Host
}
