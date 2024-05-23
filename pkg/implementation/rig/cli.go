package rig

import (
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/spf13/cobra"
)

func (rpf rigHostPluginFactory) CliBuild(cmd *cobra.Command, c *host.HostsComponent) error {

	exec.CliBuild(cmd, c)

	return nil
}
