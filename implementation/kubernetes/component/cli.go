package kubernetes

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/spf13/cobra"
)

func (c *Component) CliBuild(cmd *cobra.Command, _ *project.Project) error {
	g := &cobra.Group{
		ID:    c.Name(),
		Title: c.Name(),
	}
	cmd.AddGroup(g)

	ksv := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:version", c.Name()),
		Short:   "Get the kubernetes server version for the cluster",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ki, kierr := c.kubernetesImplementation(ctx)
			if kierr != nil {
				return kierr
			}

			ksv, ksverr := ki.ServerVersion()
			if ksverr != nil {
				return ksverr
			}

			fmt.Printf("%+v", ksv)
			return nil
		},
	}
	cmd.AddCommand(ksv)

	kcc := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:kubeconfig", c.Name()),
		Short:   "Get the kubeconfig for the cluster (BROKEN)",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cb := c.config.KubeConfig()
			if len(cb) == 0 {
				return fmt.Errorf("failed to build kubeconfig")
			}
			fmt.Println(string(cb))

			return nil
		},
	}
	cmd.AddCommand(kcc)

	return nil
}
