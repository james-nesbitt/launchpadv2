package exec

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/spf13/cobra"
)

func CliBuild(cmd *cobra.Command, c *host.HostsComponent) error {
	var exec_hn string
	var exec_sudo bool

	exec := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:execute", c.Name()),
		Short:   "Execute a command on a host",
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

			he := HostGetExecutor(h)
			if he == nil {
				return fmt.Errorf("%s: host is not an executor", h.Id())
			}

			if err := he.Connect(ctx); err != nil {
				return fmt.Errorf("%s: host failed to connect", h.Id())
			}

			opts := ExecOptions{
				Sudo: exec_sudo,
			}

			o, e, err := he.Exec(ctx, strings.Join(args, " "), os.Stdin, opts)
			if err != nil {
				return fmt.Errorf("%s: exec error: %s :: %s", h.Id(), err.Error(), e)
			}
			fmt.Println(o)

			return nil
		},
	}
	exec.Flags().StringVar(&exec_hn, "host", "", "host to execute on")
	exec.Flags().BoolVar(&exec_sudo, "sudo", false, "run command with sudo")
	cmd.AddCommand(exec)

	execi := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:interactive", c.Name()),
		Short:   "Execute interactive shell on a host",
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

			he := HostGetExecutor(h)
			if he == nil {
				return fmt.Errorf("%s: host is not an executor", h.Id())
			}

			if err := he.Connect(ctx); err != nil {
				return fmt.Errorf("%s: host failed to connect", h.Id())
			}

			opts := ExecOptions{}

			err := he.ExecInteractive(ctx, opts)
			if err != nil {
				return fmt.Errorf("%s: exec error: %s", h.Id(), err.Error())
			}

			return nil
		},
	}
	execi.Flags().StringVar(&exec_hn, "host", "", "host to execute on")
	cmd.AddCommand(execi)

	return nil
}
