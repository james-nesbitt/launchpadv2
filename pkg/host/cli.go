package host

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func (c HostsComponent) CliBuild(cmd *cobra.Command) error {

	g := &cobra.Group{
		ID:    c.Name(),
		Title: "Hosts",
	}
	cmd.AddGroup(g)

	var exec_hn string
	exec := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:execute", c.Name()),
		Short:   "Execute a command on a host",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			h := c.hosts.Get(exec_hn)
			if h == nil {

				for id, h := range c.hosts {
					slog.Info("host", slog.String("id", id), slog.Any("host", h))
				}
				return fmt.Errorf("%s: host not found", exec_hn)
			}

			o, e, err := h.Exec(ctx, strings.Join(args, " "), os.Stdin, ExecOptions{})
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
