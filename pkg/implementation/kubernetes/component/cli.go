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

			ksv, ksverr := ki.ServerVersion(ctx)
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
		Short:   "Get the kubeconfig for the cluster",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("TODO: kubernetes component does not yet have a way to generate kubeconf (we don't have config defined yet)")
		},
	}
	cmd.AddCommand(kcc)

	return nil
}
