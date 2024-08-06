package host

import (
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/spf13/cobra"
)

type HostFactoryCliBuilder interface {
	CliBuild(cmd *cobra.Command, c *HostsComponent) error
}

// CliBuild build command cli for launchpad
//
// The hosts component allows HostPluginFactories to also build CLI.
func (c *HostsComponent) CliBuild(cmd *cobra.Command, _ *project.Project) error {

	g := &cobra.Group{
		ID:    c.Name(),
		Title: "Hosts",
	}
	cmd.AddGroup(g)

	for key, hpf := range hostPluginFactories {
		hpfcb, ok := hpf.(HostFactoryCliBuilder)
		if !ok {
			slog.Debug(fmt.Sprintf("%s host plugin factory doesn't build clis", key))
			continue
		}
		if hpfcb == nil { // this should never happen
			slog.Error(fmt.Sprintf("%s host plugin factory cli builder is nil", key))
			continue
		}

		slog.Debug(fmt.Sprintf("%s host plugin factory builds clis", key))
		if err := hpfcb.CliBuild(cmd, c); err != nil {
			slog.Warn(fmt.Sprintf("%s host plugin factory cli build error", key))
		}
	}

	return nil
}
