/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/action"
)

// statusCmd represents the status command.
var statusCmd = &cobra.Command{
	GroupID: "project",
	Use:     "project:status",
	Short:   "Evaluate the component status for your project",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if valerr := cl.Validate(ctx); valerr != nil {
			return fmt.Errorf("project validation error: %s", valerr.Error())
		}

		sc, err := cl.Command(ctx, action.CommandKeyDiscover)
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

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
