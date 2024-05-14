/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Mirantis/launchpad/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the apply command.
var versionCmd = &cobra.Command{
	GroupID: "about",
	Use:     "about:version",
	Short:   "Print version information about launchpad (json)",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _ := json.Marshal(struct {
			Version     string `json:"version"`
			GitCommit   string `json:"git_commit"`
			Environment string `json:"environment"`
			GitHubRepo  string `json:"repo"`
		}{
			Version:     version.Version,
			GitCommit:   version.GitCommit,
			Environment: version.Environment,
			GitHubRepo:  version.GitHubRepo,
		})

		fmt.Println(string(v))
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
