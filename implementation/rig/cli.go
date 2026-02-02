package rig

import (
	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

// CliBuild build the RIG host plugin cli coimmands.
func (_ HostPluginFactory) CliBuild(cmd *cobra.Command, c *host.HostsComponent) error { //nolint:staticcheck // needed for insterface
	exec.CliBuild(cmd, c)

	return nil
}
