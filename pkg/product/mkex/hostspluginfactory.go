package mkex

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/product/mkex/implementation"
	"github.com/spf13/cobra"
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMKEX, &hostPluginFactory{})
}

type hostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin.
func (pf *hostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &hostPlugin{
		h: h,
	}
	pf.ps = append(pf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin.
func (pf *hostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &hostPlugin{h: h}
	pf.ps = append(pf.ps, p)

	err := d(p)

	return p, err
}

func (pf *hostPluginFactory) CliBuild(cmd *cobra.Command, c *host.HostsComponent) error {
	var exec_hn string
	exec := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:rpm-ostree:status", c.Name()),
		Short:   "Retrieve host rpm-ostree status",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			h := c.Hosts().Get(exec_hn)
			if h == nil {
				for id, h := range c.Hosts() {
					slog.Info("host", slog.String("id", id), slog.Any("host", h))
				}
				return fmt.Errorf("%s: host not found", exec_hn)
			}

			he := exec.HostGetExecutor(h)
			rosi := implementation.NewRpmOSTreeExecClient(he, exec.ExecOptions{})

			s, serr := rosi.Status(ctx)
			if serr != nil {
				return fmt.Errorf("%s: rpm-ostree status error: %s", h.Id(), serr.Error())
			}

			so, _ := json.Marshal(s)
			fmt.Println(so)

			return nil
		},
	}
	exec.Flags().StringVar(&exec_hn, "host", "", "host to execute on")
	cmd.AddCommand(exec)

	return nil
}
