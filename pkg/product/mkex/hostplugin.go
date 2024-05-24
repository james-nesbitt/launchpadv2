package mkex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/product/mkex/implementation"
	"github.com/spf13/cobra"
)

/**
 * MCR has a HostPlugin which allows MCR specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleMKEX = "mkex"
)

var (
	// ManagerHostRoles the Host roles accepted for managers.
	ManagerHostRoles = []string{"manager"}
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMKEX, &hostPluginFactory{})
}

type hostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin
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
// yaml/json .HostPluginDecode() function, and turn it into a plugin
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

// Get the MKEX plugin from a Host
func HostGetMKEX(h *host.Host) *hostPlugin {
	hgmcr := h.MatchPlugin(HostRoleMKEX)
	if hgmcr == nil {
		return nil
	}

	hmcr, ok := hgmcr.(*hostPlugin)
	if !ok {
		return nil
	}

	return hmcr
}

// hostPlugin
//
// Implements no generic plugin interfaces.
type hostPlugin struct {
	h          *host.Host
	SwarmRole  string `yaml:"role"`
	DaemonJson string `yaml:"daemon_json"`
}

func (_ hostPlugin) Id() string {
	return "mcr"
}

func (_ hostPlugin) Validate() error {
	return nil
}

func (hc hostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMKEX:
		return true
	case dockerhost.HostRoleDockerExec:
		// pkg/implementation/docker/host
		return true
	}

	return false
}

// DockerExecOptions meet the dockerhost.HostDockerExec interface
func (hc hostPlugin) DockerExec() *dockerimplementation.DockerExec {
	e := exec.HostGetExecutor(hc.h)
	if e == nil {
		return nil
	}

	def := func(ctx context.Context, cmd string, i io.Reader, rops dockerimplementation.RunOptions) (string, string, error) {
		hopts := exec.ExecOptions{Sudo: false}

		if rops.ShowOutput {
			hopts.OutputLevel = "info"
		}
		if rops.ShowError {
			hopts.ErrorLevel = "warn"
		}

		return e.Exec(ctx, cmd, i, hopts)
	}

	return dockerimplementation.NewDockerExec(def)
}

func (hc hostPlugin) DockerConfig() string {
	return hc.DaemonJson
}

func (hc hostPlugin) IsManager() bool {
	for _, r := range ManagerHostRoles {
		if hc.SwarmRole == r {
			return true
		}
	}
	return false
}
