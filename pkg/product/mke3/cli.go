package mke3

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	dockertypessystem "github.com/docker/docker/api/types/system"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

const (
	DefaultLeaderMananger = "<default>"
)

func (c MKE3) CliBuild(cmd *cobra.Command) error {

	g := &cobra.Group{
		ID:    c.Name(),
		Title: c.Name(),
	}
	cmd.AddGroup(g)

	mked := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:discover", c.Name()),
		Short:   "Discover MKE state",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			mhs, gherr := c.GetManagerHosts(ctx)
			if gherr != nil {
				return fmt.Errorf("MCR has no hosts to discover: %s", gherr.Error())
			}

			fmt.Println(fmt.Sprintf("%+v", mhs))

			return nil
		},
	}
	cmd.AddCommand(mked)

	var install_mn string
	mkei := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:install", c.Name()),
		Short:   "Install MKE3",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			slog.InfoContext(ctx, "MKE3 install starting. Looking for a manager to operate on")

			mhs, gherr := c.GetManagerHosts(ctx)
			if gherr != nil {
				return fmt.Errorf("MCR has no hosts to discover: %s", gherr.Error())
			}

			var m *dockerhost.Host
			var mi dockertypessystem.Info

			if install_mn != DefaultLeaderMananger {
				slog.InfoContext(ctx, fmt.Sprintf("looking for host %s as instructed", install_mn))
				h := mhs.Get(install_mn)
				if h == nil {
					return fmt.Errorf("%s: host not found", install_mn)
				}
				i, ierr := h.Docker(ctx).Info(ctx)
				if ierr != nil {
					return fmt.Errorf("%s: host is not a docker machine", install_mn)
				}

				m = h
				mi = i
			} else {
				slog.InfoContext(ctx, "no suggested manager, looking for first swarm manager", slog.Any("managers", mhs))
				for _, h := range mhs {
					i, ierr := h.Docker(ctx).Info(ctx)
					if ierr != nil {
						slog.WarnContext(ctx, fmt.Sprintf("%s: host is not a docker machine", h.Id()), slog.Any("host", h))
						continue
					}

					if i.Swarm.ControlAvailable {
						slog.DebugContext(ctx, fmt.Sprintf("%s: this host can act as a leader.", h.Id()), slog.Any("host", h))

						m = h
						mi = i

						break
					}
					slog.WarnContext(ctx, fmt.Sprintf("%s: host rejected as leader", h.Id()), slog.Any("info", i))
				}
			}

			if m == nil {
				return fmt.Errorf("no swarm leader found")
			}

			slog.DebugContext(ctx, fmt.Sprintf("%s: Leader found for installation", m.Id()), slog.Any("info", mi))

			if err := mkeInstall(ctx, m, c.config); err != nil {
				return fmt.Errorf("%s: bootstrap fail: %s", m.Id(), err.Error())
			}

			slog.InfoContext(ctx, "MKE3 install completed")
			return nil
		},
	}
	mkei.Flags().StringVar(&install_mn, "manager", DefaultLeaderMananger, "manager host to bootstrap")
	cmd.AddCommand(mkei)

	var uninstall_mn string
	mkeu := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:uninstall", c.Name()),
		Short:   "Uninstall MKE3",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			slog.InfoContext(ctx, "MKE3 Uninstall starting. Looking for a manager to operate on")

			slog.DebugContext(ctx, "discovering cluster (retrieving docker info)")
			mhs, gherr := c.GetManagerHosts(ctx)
			if gherr != nil {
				return fmt.Errorf("MCR has no hosts to discover: %s", gherr.Error())
			}

			var m *dockerhost.Host
			var mi dockertypessystem.Info

			if uninstall_mn != DefaultLeaderMananger {
				slog.InfoContext(ctx, fmt.Sprintf("looking for host %s as instructed", uninstall_mn))
				h := mhs.Get(uninstall_mn)
				if h == nil {
					return fmt.Errorf("%s: host not found", uninstall_mn)
				}
				i, ierr := h.Docker(ctx).Info(ctx)
				if ierr != nil {
					return fmt.Errorf("%s: host is not a docker machine", uninstall_mn)
				}

				m = h
				mi = i
			} else {
				slog.InfoContext(ctx, "no suggested manager, looking for first swarm manager", slog.Any("managers", mhs))
				for _, h := range mhs {
					i, ierr := h.Docker(ctx).Info(ctx)
					if ierr != nil {
						slog.WarnContext(ctx, fmt.Sprintf("%s: host is not a docker machine", h.Id()), slog.Any("host", h))
						continue
					}

					if i.Swarm.ControlAvailable {
						slog.DebugContext(ctx, fmt.Sprintf("%s: this host can act as a leader.", h.Id()), slog.Any("host", h))

						m = h
						mi = i

						break
					}
					slog.WarnContext(ctx, fmt.Sprintf("%s: host rejected as leader", h.Id()), slog.Any("info", i))
				}
			}

			if m == nil {
				return fmt.Errorf("no swarm leader found")
			}

			slog.DebugContext(ctx, fmt.Sprintf("%s: Leader found for un-installation", m.Id()), slog.Any("info", mi))

			if err := mkeUninstall(ctx, m, c.config); err != nil {
				return fmt.Errorf("%s: uninstall fail: %s", m.Id(), err.Error())
			}

			slog.InfoContext(ctx, "Pruning nodes after MKE uninstall")
			if err := mkePruneAfterInstall(ctx, mhs); err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("system prune after uninstall failed: %s", err.Error()))
			}

			slog.InfoContext(ctx, "MKE3 Uninstall completed")
			return nil
		},
	}
	mkeu.Flags().StringVar(&uninstall_mn, "manager", DefaultLeaderMananger, "manager host to use for uninstall")
	cmd.AddCommand(mkeu)

	return nil
}
