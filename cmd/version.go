/*
Copyright © 2023 Mirantis <tools@mirantis.com>
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
