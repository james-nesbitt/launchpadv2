package rig

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host/exec"
)

// Discover if the host is available.
func (hp hostPlugin) Discover(ctx context.Context) error {
	return hp.rig.Connect(ctx)
}

// MachineID uniquely identify the machine.
func (hp hostPlugin) MachineID(ctx context.Context) (string, error) {
	o, _, err := hp.Exec(ctx, `cat /etc/machine-id || cat /var/lib/dbus/machine-id`, nil, exec.ExecOptions{})
	return o, err
}
