package rig

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/spf13/cobra"
)

func (rpf rigHostPluginFactory) CliBuild(cmd *cobra.Command, c *host.HostsComponent) error {

	var exec_hn string
	exec := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:execute", c.Name()),
		Short:   "[rig] Execute a command on a host",
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
			if he == nil {
				return fmt.Errorf("%s: host is not an executor", h.Id())
			}

			o, e, err := he.Exec(ctx, strings.Join(args, " "), os.Stdin, exec.ExecOptions{})
			if err != nil {
				return fmt.Errorf("%s: exec error: %s :: %s", h.Id(), err.Error(), e)
			}
			fmt.Println(o)

			return nil
		},
	}
	exec.Flags().StringVar(&exec_hn, "host", "", "host to execute on")
	cmd.AddCommand(exec)

	return nil
}
