/*
t
Copyright © 2023 Mirantis <tools@mirantis.com>
*/
package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/cmd"
	"github.com/Mirantis/launchpad/pkg/project"

	// Register v2 spec handler.
	_ "github.com/Mirantis/launchpad/mirantis/config/v2_0"
	_ "github.com/Mirantis/launchpad/mirantis/config/v2_1"

	// Register Host handlers.
	_ "github.com/Mirantis/launchpad/implementation/rig"

	// Register product handlers.
	_ "github.com/Mirantis/launchpad/mirantis/product/mke3"
	_ "github.com/Mirantis/launchpad/mirantis/product/mkex"
	_ "github.com/Mirantis/launchpad/mirantis/product/msr2"

	// Register mock stuff, in case it gets used.
	_ "github.com/Mirantis/launchpad/pkg/mock"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mkex",
		Short: "A Mirantis manager of MKEx cluster",
		Long:  `Manage MKEx clusters and their components..`,
	}

	slog.Debug("bootstrapping cli project from config file")
	proj := cmd.Bootstrap(rootCmd)

	slog.Debug("bootstrapping cli from projec")
	if err := project.CliBuild(rootCmd, &proj); err != nil {
		slog.Error("failed to bootstrap cli")
	}

	slog.Debug("bootstrapping cli complete")

	rootCmd.ParseFlags(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to run command")
		os.Exit(1)
	}
}
