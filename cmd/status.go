/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/config"
)

// statusCmd represents the status command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		f, foerr := os.Open(cfgFile)
		if foerr != nil {
			return fmt.Errorf("could not open config '%s' : %s", cfgFile, foerr.Error())
		}

		yb, frerr := io.ReadAll(f)
		if frerr != nil {
			return fmt.Errorf("could not read config '%s' : %s", cfgFile, frerr.Error())
		}

		cl, umerr := config.ConfigFromYamllBytes(yb)
		if umerr != nil {
			return fmt.Errorf("Error occured unarshalling yaml: %s \nYAML:\b%s", umerr.Error(), yb)
		}

		if valerr := cl.Validate(ctx); valerr != nil {
			return fmt.Errorf("cluster validation error: %s", valerr.Error())
		}

		sc, err := cl.Command(ctx, action.CommandKeyDiscover)
		if err != nil {
			return fmt.Errorf("cluster command build error: %s", err.Error())
		}

		if err := sc.Validate(ctx); err != nil {
			return fmt.Errorf("cluster command [%s] validation failed: %s", sc.Key, err.Error())
		}

		if err := sc.Run(ctx); err != nil {
			return fmt.Errorf("cluster command [%s] execution failed: %s", sc.Key, err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
