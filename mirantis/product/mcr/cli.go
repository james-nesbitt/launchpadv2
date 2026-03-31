package mcr

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/spf13/cobra"

	dockertypesswarm "github.com/docker/docker/api/types/swarm"
	dockertypessystem "github.com/docker/docker/api/types/system"

	dockerhost "github.com/Mirantis/launchpad/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/project"
)

func (c Component) CliBuild(cmd *cobra.Command, _ project.Project) error {
	g := &cobra.Group{
		ID:    c.Name(),
		Title: c.Name(),
	}
	cmd.AddGroup(g)

	mcrhsd := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:discover", c.Name()),
		Short:   "Discover MCR state",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, gherr := c.GetAllHosts(ctx)
			if gherr != nil {
				return fmt.Errorf("mcr has no hosts to discover: %s", gherr.Error())
			}

			info := map[string]dockertypessystem.Info{}
			infomu := sync.Mutex{}

			if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
				if err := exec.HostGetExecutor(h).Connect(ctx); err != nil {
					return fmt.Errorf("failed to connect to host %s: %w", h.ID(), err)
				}
				slog.InfoContext(ctx, fmt.Sprintf("%s: discovering MCR state", h.ID()), slog.Any("host", h))

				i, err := dockerhost.HostGetDockerExec(h).Info(ctx)
				if err != nil {
					slog.WarnContext(ctx, fmt.Sprintf("%s: MCR state discovery failure", h.ID()), slog.Any("host", h), slog.Any("error", err))
					return fmt.Errorf("%s: failed to update docker info: %s", h.ID(), err.Error())
				}

				infomu.Lock()
				info[h.ID()] = i
				infomu.Unlock()

				return nil
			}); err != nil {
				return fmt.Errorf("docker info update failed: %s", err.Error())
			}

			o, _ := json.Marshal(info)
			fmt.Println(string(o))

			return nil
		},
	}
	cmd.AddCommand(mcrhsd)

	mcrswd := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:swarm", c.Name()),
		Short:   "Discover MCR Swarm",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			mhs, gherr := c.GetManagerHosts(ctx)
			if gherr != nil {
				return fmt.Errorf("mcr has no hosts to discover: %s", gherr.Error())
			}

			var l *host.Host

			for _, h := range mhs {
				i, ierr := dockerhost.HostGetDockerExec(h).Info(ctx)
				if ierr != nil {
					continue
				}

				if i.Swarm.ControlAvailable {
					slog.DebugContext(ctx, fmt.Sprintf("%s: swarm already active, this host can act as a leader.", h.ID()), slog.Any("host", h))
					l = h
					break
				}
			}

			if l == nil {
				return fmt.Errorf("no swarm leader found")
			}

			ld := dockerhost.HostGetDockerExec(l)

			li, lierr := ld.Info(ctx)
			if lierr != nil {
				return fmt.Errorf("%s: swarm join failed because leader docker info error: %s", l.ID(), lierr.Error())
			}
			si, sierr := ld.SwarmInspect(ctx)
			if sierr != nil {
				return fmt.Errorf("%s: swarm join failed because leader docker swarm inspect error: %s", l.ID(), sierr.Error())
			}
			ni, nierr := ld.NodeList(ctx, dockertypesswarm.NodeListOptions{})
			if nierr != nil {
				return fmt.Errorf("%s: swarm join failed because leader docker swarm node list error: %s", l.ID(), nierr.Error())
			}

			info := struct {
				Li dockertypessystem.Info
				Si dockertypesswarm.Swarm
				Ni []dockertypesswarm.Node
			}{
				Li: li,
				Si: si,
				Ni: ni,
			}

			o, _ := json.Marshal(info)
			fmt.Println(string(o))

			return nil
		},
	}
	cmd.AddCommand(mcrswd)

	return nil
}
