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
	// Arch what CPU Architecture is the host using
	Arch(ctx context.Context) string
	// UserCacheDir returns the default root directory to use for user-specific cached data.
	UserCacheDir() string
	// UserConfigDir returns the default root directory to use for user-specific configuration data.
	UserConfigDir() string
	// UserHomeDir returns the current user's home directory.
	UserHomeDir() string
}
