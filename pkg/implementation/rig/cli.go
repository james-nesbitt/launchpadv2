package rig

import (
	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

func (_ HostPluginFactory) CliBuild(cmd *cobra.Command, c *host.HostsComponent) error {

	exec.CliBuild(cmd, c)

	return nil
}
