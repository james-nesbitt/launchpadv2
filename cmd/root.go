/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/cluster"
	"github.com/Mirantis/launchpad/pkg/config"
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
	cfgFile string
	cl      cluster.Cluster
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "A Mirantis installer",
	Long:  `Install various Mirantis products.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

		cl = tcl

		if valerr := cl.Validate(cmd.Context()); valerr != nil {
			return fmt.Errorf("cluster validation error: %s", valerr.Error())
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.launchpad.yaml)")
}
