package host

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

type HostFactoryCliBuilder interface {
	CliBuild(cmd *cobra.Command, c *HostsComponent) error
}

// CliBuild build command cli for launchpad
//
// The hosts component allows HostPluginFactories to also build CLI
func (c *HostsComponent) CliBuild(cmd *cobra.Command) error {

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

		if err := hpfcb.CliBuild(cmd, c); err != nil {
			slog.Warn(fmt.Sprintf("%s host plugin factory cli build error", key))
		}
	}

	return nil
}
