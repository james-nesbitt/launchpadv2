package rig

import "context"

// Discover if the host is available
func (hp rigHostPlugin) Discover(ctx context.Context) error {
	return hp.rig.Connect(ctx)
}
