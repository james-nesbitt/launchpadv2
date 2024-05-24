package project

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

const (
	ProjectCliKey = "project"
)

func CliBuild(cmd *cobra.Command, p *Project) error {
	slog.DebugContext(cmd.Context(), "building cli commands from project")

	g := &cobra.Group{
		ID:    ProjectCliKey,
		Title: "Project",
	}
	cmd.AddGroup(g)

	projectCliBuildStatus(cmd, p)
	projectCliBuildApply(cmd, p)
	projectCliBuildReset(cmd, p)

	if len(p.Components) == 0 {
		slog.WarnContext(cmd.Context(), "empty project object was built, so no project/component commands are available.")
		return nil
	}

	slog.DebugContext(cmd.Context(), "building cli commands from project components")
	for _, c := range p.Components {
		if ccb, ok := c.(CliBuilder); ok {
			if err := ccb.CliBuild(cmd, p); err != nil {
				slog.ErrorContext(cmd.Context(), fmt.Sprintf("%s: error when building cli", c.Name()), slog.Any("component", c))
			}
		}
	}

	return nil
}

type CliBuilder interface {
	CliBuild(*cobra.Command, *Project) error
}

func projectCliBuildStatus(cmd *cobra.Command, p *Project) error {
	// statusCmd represents the status command.
	cmd.AddCommand(&cobra.Command{
		GroupID: ProjectCliKey,
		Use:     "project:status",
		Short:   "Evaluate the component status for your project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if valerr := p.Validate(ctx); valerr != nil {
				return fmt.Errorf("project validation error: %s", valerr.Error())
			}

			sc, err := p.Command(ctx, CommandKeyDiscover)
			if err != nil {
				return fmt.Errorf("project command build error: %s", err.Error())
			}

			if err := sc.Validate(ctx); err != nil {
				return fmt.Errorf("project command [%s] validation failed: %s", sc.Key, err.Error())
			}

			if err := sc.Run(ctx); err != nil {
				return fmt.Errorf("project command [%s] execution failed: %s", sc.Key, err.Error())
			}

			return nil
		},
	})
	return nil
}

func projectCliBuildApply(cmd *cobra.Command, p *Project) error {

	// applyCmd represents the apply command.
	cmd.AddCommand(&cobra.Command{
		GroupID: ProjectCliKey,
		Use:     "project:apply",
		Short:   "Apply any component installs on your project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			sc, err := p.Command(ctx, CommandKeyApply)
			if err != nil {
				return fmt.Errorf("project command build error: %s", err.Error())
			}

			if err := sc.Validate(ctx); err != nil {
				return fmt.Errorf("project command [%s] validation failed: %s", sc.Key, err.Error())
			}

			if err := sc.Run(ctx); err != nil {
				return fmt.Errorf("project command [%s] execution failed: %s", sc.Key, err.Error())
			}

			return nil
		},
	})
	return nil
}

func projectCliBuildReset(cmd *cobra.Command, p *Project) error {
	cmd.AddCommand(&cobra.Command{
		GroupID: ProjectCliKey,
		Use:     "project:reset",
		Short:   "Remove any component installs in your project",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if valerr := p.Validate(ctx); valerr != nil {
				return fmt.Errorf("project validation error: %s", valerr.Error())
			}

			sc, err := p.Command(ctx, CommandKeyReset)
			if err != nil {
				return fmt.Errorf("project command build error: %s", err.Error())
			}

			if err := sc.Validate(ctx); err != nil {
				return fmt.Errorf("project command [%s] validation failed: %s", sc.Key, err.Error())
			}

			if err := sc.Run(ctx); err != nil {
				return fmt.Errorf("project command [%s] execution failed: %s", sc.Key, err.Error())
			}

			return nil
		},
	})
	return nil
}
