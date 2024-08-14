package msr4

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

	sc := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:dev-install", c.Name()),
		Short:   "msr4 install (used during development)",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ki, kierr := c.GetKubernetesImplementation(ctx)
			if kierr != nil {
				return kierr
			}
			fmt.Printf("Retrieved Kube implementation (from whoever provided kube): %+v \n", ki)

			hc, hcerr := ki.HelmClient(ctx, c.helmOptions())
			if hcerr != nil {
				return hcerr
			}
			fmt.Printf("Built helm client: %+v \n", hc)

			hr := c.helmRepo()

			fmt.Printf("Adding helm repo: %+v \n", hr)
			if err := hc.AddOrUpdateChartRepo(hr); err != nil {
				return fmt.Errorf("Error adding repo: %s", err.Error())
			}

			hcs := c.helmChartSpec()
			hcopts := c.helmChartOpts()

			fmt.Sprintf("Installing/Upgrading chart: chart:%+v opts:%+v \n", hcs, hcopts)
			r, rerr := hc.InstallOrUpgradeChart(ctx, &hcs, &hcopts)
			if rerr != nil {
				return rerr
			}
			fmt.Printf("Release instaled: %+v \n", *r)

			return nil
		},
	}
	cmd.AddCommand(sc)

	return nil
}
