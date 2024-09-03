package mkex

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/spf13/cobra"
)

func (c *Component) CliBuild(cmd *cobra.Command, _ *project.Project) error {

	g := &cobra.Group{
		ID:    c.Name(),
		Title: c.Name(),
	}
	cmd.AddGroup(g)

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:discover", c.Name()),
		Short:   "Discover MKEX state",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, gherr := c.GetAllHosts(ctx)
			if gherr != nil {
				return fmt.Errorf("MKEX has no hosts to discover: %s", gherr.Error())
			}

			fmt.Println(fmt.Sprint(hs))

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:activate-swarm", c.Name()),
		Short:   "Actiate Swarm on the MKEX cluster",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return activateSwarm(cmd.Context(), c)
		},
	})

	return nil
}
