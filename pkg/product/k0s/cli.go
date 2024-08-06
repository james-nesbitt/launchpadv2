package k0s

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func (c *Component) CliBuild(cmd *cobra.Command, _ *project.Project) error {
	var hn string // host name
	var ln string // leader host name
	var r string  // host role

	g := &cobra.Group{
		ID:    c.Name(),
		Title: c.Name(),
	}
	cmd.AddGroup(g)

	sc := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:status", c.Name()),
		Short:   "k0s status",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, hserr := c.GetAllHosts(ctx)
			if hserr != nil {
				return hserr
			}

			h := hs.Get(hn)
			if h == nil {
				return fmt.Errorf("could not find host %s", hn)
			}

			kh := HostGetK0s(h)

			s, err := kh.Status(ctx)
			if err != nil {
				if errors.Is(err, ErrK0sBinaryNotFound) {
					return ErrK0sBinaryNotFound
				}
				if errors.Is(err, ErrK0sNotRunning) {
					return ErrK0sNotRunning
				}
				return err
			}

			y, _ := yaml.Marshal(s)
			fmt.Println(string(y))
			return nil
		},
	}
	sc.Flags().StringVar(&hn, "host", "", "host to execute on")
	cmd.AddCommand(sc)

	vc := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:version", c.Name()),
		Short:   "k0s version",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, hserr := c.GetAllHosts(ctx)
			if hserr != nil {
				return hserr
			}
			h := hs.Get(hn)
			if h == nil {
				return fmt.Errorf("could not find host %s", hn)
			}

			kh := HostGetK0s(h)
			if kh == nil {
				return fmt.Errorf("host %s is not a k0s host", hn)
			}

			v, err := kh.Version(ctx)
			if err != nil {
				if errors.Is(err, ErrK0sBinaryNotFound) {
					return ErrK0sBinaryNotFound
				}
				return err
			}

			j, _ := json.Marshal(v)
			fmt.Println(string(j))
			return nil
		},
	}
	vc.Flags().StringVar(&hn, "host", "", "host to execute on")
	cmd.AddCommand(vc)

	cmd.AddCommand(&cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:get-k0s", c.Name()),
		Short:   "get the k0s binary on to each host",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, hserr := c.GetAllHosts(ctx)
			if hserr != nil {
				return hserr
			}

			if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
				kh := HostGetK0s(h)

				if c.config.ShouldDownload() {
					slog.InfoContext(ctx, fmt.Sprintf("%s: downloading binary to host", h.Id()))
					if err := kh.DownloadK0sBinary(ctx, c.config.Version); err != nil {
						return err
					}
				} else {
					slog.InfoContext(ctx, fmt.Sprintf("%s: uploading binary to host", h.Id()))
					if err := kh.UploadK0sBinary(ctx, c.config.Version); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				return nil
			}

			return nil
		},
	})

	cc := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:config", c.Name()),
		Short:   "build host config",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, hserr := c.GetAllHosts(ctx)
			if hserr != nil {
				return hserr
			}
			h := hs.Get(hn)
			if h == nil {
				return fmt.Errorf("%s: host not found", hn)
			}

			baseCfg := c.config.K0sConfig
			csans := c.CollectClusterSans(ctx)

			kh := HostGetK0s(h)

			cfg, cerr := kh.BuildHostConfig(ctx, baseCfg, csans)
			if cerr != nil {
				return cerr
			}

			cfgbs, _ := yaml.Marshal(cfg)
			fmt.Println(string(cfgbs))
			return nil
		},
	}
	cc.Flags().StringVar(&hn, "host", "", "host to execute on")
	cmd.AddCommand(cc)

	ac := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:activate", c.Name()),
		Short:   "activate k0s cluster",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hs, hserr := c.GetControllerHosts(ctx)
			if hserr != nil {
				return hserr
			}
			l := hs.Get(ln)
			if l == nil {
				return fmt.Errorf("%s: leader host not found", ln)
			}
			slog.InfoContext(ctx, fmt.Sprintf("%s: using as leader", l.Id()))

			baseCfg := c.config.K0sConfig
			csans := c.CollectClusterSans(ctx)

			lkh := HostGetK0s(l)

			slog.InfoContext(ctx, fmt.Sprintf("%s: writing config to leader host", l.Id()))
			if werr := lkh.WriteK0sConfig(ctx, baseCfg, csans); werr != nil {
				return werr
			}

			slog.InfoContext(ctx, fmt.Sprintf("%s: activating leader host", l.Id()))
			if err := lkh.ActivateNewCluster(ctx, c.config); err != nil {
				return err
			}

			fmt.Printf("%s: started new cluster", l.Id())
			return nil
		},
	}
	ac.Flags().StringVar(&ln, "leader", "", "host to activate as leader")
	cmd.AddCommand(ac)

	jc := &cobra.Command{
		GroupID: c.Name(),
		Use:     fmt.Sprintf("%s:join", c.Name()),
		Short:   "join host to k0s cluster",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var l *host.Host
			var hs host.Hosts
			var hserr error

			if ln == "" {
				l = c.GetLeaderHost(ctx)
			} else {
				hs, hserr = c.GetControllerHosts(ctx)
				if hserr != nil {
					return fmt.Errorf("could not retrieve any controllers")
				}
				l = hs.Get(ln)
			}
			if l == nil {
				return fmt.Errorf("could not find a cluster controller/leader")
			}

			switch r {
			case RoleController:
				// we already collected hosts
			case RoleWorker:
				hs, hserr = c.GetWorkerHosts(ctx)
			default:
				return fmt.Errorf("unrecognized role for join: %s", r)
			}

			if hserr != nil {
				return hserr
			}
			h := hs.Get(hn)
			if h == nil {
				return fmt.Errorf("%s: host not found (%s)", hn, r)
			}

			kh := HostGetK0s(h)

			switch r {
			case RoleController:
				baseCfg := c.config.K0sConfig
				csans := c.CollectClusterSans(ctx)

				slog.InfoContext(ctx, fmt.Sprintf("%s: writing config to controller host", h.Id()))
				if werr := kh.WriteK0sConfig(ctx, baseCfg, csans); werr != nil {
					return werr
				}

			case RoleWorker:
			}

			slog.InfoContext(ctx, fmt.Sprintf("%s: joining host as %s to leader %s", h.Id(), r, l.Id()))
			if err := kh.JoinCluster(ctx, l, r, c.config); err != nil {
				return err
			}

			fmt.Printf("%s: joined cluster using leader %s", h.Id(), l.Id())
			return nil
		},
	}
	jc.Flags().StringVar(&ln, "leader", "", "host to use as leader")
	jc.Flags().StringVar(&hn, "host", "", "host to join to cluster")
	jc.Flags().StringVar(&r, "role", "worker", "role to join to cluster")
	cmd.AddCommand(jc)

	return nil
}
