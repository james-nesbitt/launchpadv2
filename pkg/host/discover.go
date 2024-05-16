package host

import "context"

/**
 * Discovering a host means checking to see if it
 * is available.
 *
 * We rely on a HostPlugin perform this.
 */

const (
	HostRoleDiscover = "discover"
)

func HostGetDiscover(h *Host) HostDiscover {
	hgd := h.MatchPlugin(HostRoleDiscover)
	if hgd == nil {
		return nil
	}

	hd, ok := hgd.(HostDiscover)
	if !ok {
		return nil
	}

	return hd
}

// HostDiscover validate that the Host is available
type HostDiscover interface {
	Discover(context.Context) error
}
