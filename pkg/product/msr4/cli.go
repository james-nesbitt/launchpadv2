package msr4

import (
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/spf13/cobra"
)

func (c *Component) CliBuild(cmd *cobra.Command, _ *project.Project) error {
	g := &cobra.Group{
		ID:    c.Name(),
		Title: c.Name(),
	}
	cmd.AddGroup(g)

	dic := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:dev-install", c.Name()),
		Short:   "msr4 install (used during development)",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			slog.InfoContext(ctx, "Retrieving Kube implementation\n")
			ki, kierr := c.GetKubernetesImplementation(ctx)
			if kierr != nil {
				return kierr
			}
			slog.InfoContext(ctx, fmt.Sprintf("Retrieved Kube implementation (from whoever provided kube): %+v \n", ki))

			hc, hcerr := ki.HelmClient(ctx, c.helmOptions())
			if hcerr != nil {
				return hcerr
			}
			slog.InfoContext(ctx, fmt.Sprintf("Built helm client: %+v \n", hc))

			hr := c.helmRepo()

			slog.InfoContext(ctx, fmt.Sprintf("Adding helm repo: %+v \n", hr))
			if err := hc.AddOrUpdateChartRepo(hr); err != nil {
				return fmt.Errorf("Error adding repo: %s", err.Error())
			}

			hcs := c.helmChartSpec()
			hcopts := c.helmChartOpts()

			slog.InfoContext(ctx, fmt.Sprintf("Installing/Upgrading chart: chart:%+v opts:%+v \n", hcs, hcopts))
			r, rerr := hc.InstallOrUpgradeChart(ctx, &hcs, &hcopts)
			if rerr != nil {
				return rerr
			}
			slog.InfoContext(ctx, fmt.Sprintf("Release installed: %+v \n", *r))

			return nil
		},
	}
	cmd.AddCommand(dic)

	duc := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:dev-uninstall", c.Name()),
		Short:   "msr4 uninstall (used during development)",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			slog.InfoContext(ctx, "Retrieving Kube implementation\n")
			ki, kierr := c.GetKubernetesImplementation(ctx)
			if kierr != nil {
				return kierr
			}
			slog.InfoContext(ctx, fmt.Sprintf("Retrieved Kube implementation (from whoever provided kube): %+v \n", ki))

			hc, hcerr := ki.HelmClient(ctx, c.helmOptions())
			if hcerr != nil {
				return hcerr
			}
			slog.InfoContext(ctx, fmt.Sprintf("Built helm client: %+v \n", hc))

			hcs := c.helmChartSpec()
			hcopts := c.helmChartOpts()

			slog.InfoContext(ctx, fmt.Sprintf("Uninstalling chart: chart:%+v opts:%+v \n", hcs, hcopts))
			rerr := hc.UninstallRelease(&hcs)
			if rerr != nil {
				return rerr
			}
			slog.InfoContext(ctx, "Release uninstalled \n")

			return nil
		},
	}
	cmd.AddCommand(duc)

	return nil
}
