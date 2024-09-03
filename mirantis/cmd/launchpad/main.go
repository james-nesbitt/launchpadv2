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

	// register implementation components.
	_ "github.com/Mirantis/launchpad/implementation/kubernetes/component"

	// Register legacy product handlers.
	_ "github.com/Mirantis/launchpad/mirantis/product/k0s"
	_ "github.com/Mirantis/launchpad/mirantis/product/mcr"
	_ "github.com/Mirantis/launchpad/mirantis/product/mke3"
	_ "github.com/Mirantis/launchpad/mirantis/product/mke4"
	_ "github.com/Mirantis/launchpad/mirantis/product/mkex"
	_ "github.com/Mirantis/launchpad/mirantis/product/msr2"
	_ "github.com/Mirantis/launchpad/mirantis/product/msr3"
	_ "github.com/Mirantis/launchpad/mirantis/product/msr4"

	// Register some more components.
	_ "github.com/Mirantis/launchpad/implementation/kubernetes/component"

	// Register mock stuff, in case it gets used.
	_ "github.com/Mirantis/launchpad/pkg/mock"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "launchpad",
		Short: "A Mirantis installer",
		Long:  `manage various Mirantis products.`,
	}

	proj := cmd.Bootstrap(rootCmd)

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
