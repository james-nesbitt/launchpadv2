package kubernetes

import "github.com/Mirantis/launchpad/pkg/host"

// State.
type State struct {
	Leader *host.Host
}
