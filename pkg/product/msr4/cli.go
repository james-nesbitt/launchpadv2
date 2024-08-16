package msr4

import (
	"encoding/json"
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

	rdc := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:status", c.Name()),
		Short:   "msr4 release status",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			r, rerr := c.getRelease(ctx)
			if rerr != nil {
				slog.ErrorContext(ctx, "failed to add MSR4 helm release")
				return rerr
			}

			slog.InfoContext(ctx, "Added MSR4 release status retrieve")

			rbs, _ := json.Marshal(r)
			fmt.Print(string(rbs))

			return nil
		},
	}
	cmd.AddCommand(rdc)

	rac := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:add-repo", c.Name()),
		Short:   "msr4 repo add",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := c.addMsr4HelmRepo(ctx); err != nil {
				slog.ErrorContext(ctx, "failed to add MSR4 helm repo")
				return err
			}

			slog.InfoContext(ctx, "Added MSR4 Helm repo")

			return nil
		},
	}
	cmd.AddCommand(rac)

	ric := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:install-release", c.Name()),
		Short:   "msr4 chart release install",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			r, rerr := c.installMsr4Release(ctx)
			if rerr != nil {
				slog.ErrorContext(ctx, "failed to install MSR4 helm chart release")
				return rerr
			}

			slog.InfoContext(ctx, "Installed MSR4 helm chart as release", slog.Any("release-info", r.Info))

			return nil
		},
	}
	cmd.AddCommand(ric)

	duc := &cobra.Command{
		GroupID: g.ID,
		Use:     fmt.Sprintf("%s:uninstall", c.Name()),
		Short:   "msr4 uninstall release",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := c.uninstallMsr4Release(ctx); err != nil {
				slog.ErrorContext(ctx, "Failed to uninstall MSR4 helm release")
				return err
			}

			slog.InfoContext(ctx, "Uninstalled MSR4 helm release (namespace may remain.)")

			return nil
		},
	}
	cmd.AddCommand(duc)

	return nil
}
