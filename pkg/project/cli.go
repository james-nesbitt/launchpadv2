package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ProjectCliKey = "project"
)

func CliBuild(cmd *cobra.Command, p *Project) error {

	g := &cobra.Group{
		ID:    ProjectCliKey,
		Title: "Project",
	}
	cmd.AddGroup(g)

	// applyCmd represents the apply command.
	applyCmd := &cobra.Command{
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
	}
	cmd.AddCommand(applyCmd)

	resetCmd := &cobra.Command{
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
	}
	cmd.AddCommand(resetCmd)

	// statusCmd represents the status command.
	statusCmd := &cobra.Command{
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
	}
	cmd.AddCommand(statusCmd)

	return nil
}
