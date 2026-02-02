package k0s

import "github.com/Mirantis/launchpad/pkg/host"

// State information for an instance of the component.
type State struct {
	Leader *host.Host
}
