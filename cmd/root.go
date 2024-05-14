/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/config"
	"github.com/Mirantis/launchpad/pkg/project"

	// Register v2 spec handler.
	_ "github.com/Mirantis/launchpad/pkg/config/v2_0"
	// Register Host handlers.
	_ "github.com/Mirantis/launchpad/pkg/host/mock"
	_ "github.com/Mirantis/launchpad/pkg/host/rig"

	// Register legacy product handlers for testing.
	_ "github.com/Mirantis/launchpad/pkg/product/mcr"
	_ "github.com/Mirantis/launchpad/pkg/product/mke3"
	_ "github.com/Mirantis/launchpad/pkg/product/msr2"

	// Register nextgen product handlers for testing.
	_ "github.com/Mirantis/launchpad/pkg/product/k0s"
	_ "github.com/Mirantis/launchpad/pkg/product/mke4"
	_ "github.com/Mirantis/launchpad/pkg/product/msr4"
)

var (
	debug   bool
	cfgFile string
	cl      *project.Project
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "launchpad",
		Short: "A Mirantis installer",
		Long:  `Install various Mirantis products.`,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "increase logging verbosity")

	// This should early pre-populate the above flags, which we will use to build more commands. This will
	// be executed again with rootCmd.Execute()
	rootCmd.ParseFlags(os.Args[1:])

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.DebugContext(rootCmd.Context(), "building project for cli")
	if cfgFile == "" {
		slog.Warn("No project config specified. Using an empty project")
		cl = &project.Project{
			Components: component.Components{},
		}
	} else if err := boostrapBuildProject(rootCmd.Context()); err != nil {
		slog.Error("failed to build project", slog.Any("error", err))
		os.Exit(1)
	}
	slog.DebugContext(rootCmd.Context(), "finished building project for cli")

	if err := bootstrapLaunchpadCmd(rootCmd); err != nil {
		slog.Error("failed to bootstrap cli")
	}

	rootCmd.ParseFlags(os.Args[1:])

	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to run command")
		os.Exit(1)
	}
}

func boostrapBuildProject(ctx context.Context) error {
	if cfgFile == "" {
		return fmt.Errorf("no config file defined")
	}

	f, foerr := os.Open(cfgFile)
	if foerr != nil {
		return fmt.Errorf("could not access config '%s' : %s", cfgFile, foerr.Error())
	}

	yb, frerr := io.ReadAll(f)
	if frerr != nil {
		return fmt.Errorf("could not read config '%s' : %s", cfgFile, frerr.Error())
	}

	tcl, umerr := config.ConfigFromYamllBytes(yb)
	if umerr != nil {
		return fmt.Errorf("Error occurred unarshalling yaml: %s \nYAML:\b%s", umerr.Error(), yb)
	}

	cl = &tcl

	if valerr := cl.Validate(ctx); valerr != nil {
		return fmt.Errorf("project validation error: %s", valerr.Error())
	}

	return nil
}

func bootstrapLaunchpadCmd(cmd *cobra.Command) error {
	cmd.AddGroup(&cobra.Group{
		ID:    "about",
		Title: "About Launchpad",
	})

	cmd.AddCommand(versionCmd)

	if cl == nil {
		slog.WarnContext(cmd.Context(), "no project object was built")
	} else if len(cl.Components) == 0 {
		slog.WarnContext(cmd.Context(), "empty project object was built, so no project/component commands are available.")
	} else {
		cmd.AddGroup(&cobra.Group{
			ID:    "project",
			Title: "Project",
		})

		cmd.AddCommand(statusCmd)
		cmd.AddCommand(applyCmd)
		cmd.AddCommand(resetCmd)

		slog.DebugContext(cmd.Context(), "building cli commands from project")

		for _, c := range cl.Components {
			if ccb, ok := c.(CliBuilder); ok {
				if err := ccb.CliBuild(cmd); err != nil {
					slog.ErrorContext(cmd.Context(), fmt.Sprintf("%s: error when building cli", c.Name()), slog.Any("component", c))
				}
			}
		}
	}

	return nil
}

type CliBuilder interface {
	CliBuild(*cobra.Command) error
}
