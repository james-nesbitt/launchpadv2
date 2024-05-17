package rig

import "context"

// IsWindows.
func (h rigHostPlugin) IsWindows(ctx context.Context) bool {
	return h.rig.IsWindows()
}
