package mke4

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

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
		Use:     fmt.Sprintf("%s:status", c.Name()),
		Short:   "MKE4 status",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:version", c.Name()),
		Short:   "mk4 version",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:blueprint-show", c.Name()),
		Short:   "Show the MKE4 blueprint",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(c.config.Blueprint)

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:operator-yaml", c.Name()),
		Short:   "Show the MKE4 operator yaml",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			oryrb, oryrberr := operatorYaml()
			if oryrberr != nil {
				return oryrberr
			}

			fmt.Println(string(oryrb))
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:operator-install", c.Name()),
		Short:   "install the MKE4 operator",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			slog.InfoContext(ctx, "retrieving kubernetes context")
			k, kerr := c.GetKubernetesImplementation(ctx)
			if kerr != nil {
				return kerr
			}
			slog.DebugContext(ctx, "retrieved kubernetes context", slog.Any("k", k))

			slog.InfoContext(ctx, "Info installing MKE4 Operator")
			if err := OperatorInstall(ctx, *k, c.config); err != nil {
				return err
			}

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:operator-uninstall", c.Name()),
		Short:   "uninstall the MKE4 operator",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			slog.InfoContext(ctx, "retrieving kubernetes context")
			k, kerr := c.GetKubernetesImplementation(ctx)
			if kerr != nil {
				return kerr
			}
			slog.DebugContext(ctx, "retrieved kubernetes context", slog.Any("k", k))

			slog.InfoContext(ctx, "Info un-installing MKE4 Operator")
			if err := OperatorUninstall(ctx, *k, c.config); err != nil {
				return err
			}

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:blueprint", c.Name()),
		Short:   "Show the MKE4 blueprint yaml",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			bp, bperr := c.config.BlueprintResource()
			if bperr != nil {
				return bperr
			}

			bpy, byperr := yaml.Marshal(bp.Object)
			if byperr != nil {
				return byperr
			}

			fmt.Println(string(bpy))
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:install", c.Name()),
		Short:   "install MKE4",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			slog.InfoContext(ctx, "retrieving kubernetes context")
			k, kerr := c.GetKubernetesImplementation(ctx)
			if kerr != nil {
				return kerr
			}
			slog.DebugContext(ctx, "retrieved kubernetes context", slog.Any("k", k))

			slog.InfoContext(ctx, "Activating MKE4 instance", slog.String("name", c.config.Name))
			if err := ActivateOrUpdate(ctx, *k, c.config); err != nil {
				return err
			}

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:uninstall", c.Name()),
		Short:   "uninstall MKE4",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			slog.InfoContext(ctx, "retrieving kubernetes context")
			k, kerr := c.GetKubernetesImplementation(ctx)
			if kerr != nil {
				return kerr
			}
			slog.DebugContext(ctx, "retrieved kubernetes context", slog.Any("k", k))

			slog.InfoContext(ctx, "Activating MKE4 instance", slog.String("name", c.config.Name))
			if err := Deactivate(ctx, *k, c.config); err != nil {
				return err
			}

			return nil
		},
	})

	return nil
}
