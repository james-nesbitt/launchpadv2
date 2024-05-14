/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/action"
)

// applyCmd represents the apply command.
var applyCmd = &cobra.Command{
	GroupID: "project",
	Use:     "project:apply",
	Short:   "Apply any component installs on your project",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		sc, err := cl.Command(ctx, action.CommandKeyApply)
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
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
