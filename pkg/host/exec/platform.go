package exec

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

// -- PLATFORM PLUGIN

const (
	HostRolePlatform = "platform"
)

func HostGetPlatform(h *host.Host) HostPlatform {
	hgp := h.MatchPlugin(HostRolePlatform)
	if hgp == nil {
		return nil
	}

	hp, ok := hgp.(HostPlatform)
	if !ok {
		return nil
	}

	return hp
}

type HostPlatform interface {
	// IsWindows check, could require a connection
	IsWindows(context.Context) bool
}