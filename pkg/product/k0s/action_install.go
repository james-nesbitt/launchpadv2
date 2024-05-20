package k0s

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
	"github.com/Mirantis/launchpad/pkg/util/download"
	"github.com/Mirantis/launchpad/pkg/util/retry"
)

type installK0sStep struct {
	baseStep
	id    string
	d     *download.QueueDownload
	Force bool
}

func (s installK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-install", s.id)
}

func (s installK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s install step", slog.String("ID", s.id))

	s.d = download.NewQueueDownload(nil)

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return fmt.Errorf("K0s install step could not retrieve host list: %s", hserr.Error())
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: downloading binaries", h.Id()))
		if err := s.uploadBinaries(ctx, h); err != nil {
			return err
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: configuring k0s", h.Id()))
		if err := s.configureK0s(ctx, h); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Error installing k0s: %s", err.Error())
	}
	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: initialize k0s", h.Id()))
		if err := s.initializeK0s(ctx, h); err != nil {
			return err
		}

		slog.InfoContext(ctx, fmt.Sprintf("%s: starting k0s", h.Id()))
		if err := s.startK0s(ctx, h); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Error installing k0s: %s", err.Error())
	}

	return nil
}

func (s installK0sStep) uploadBinaries(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: uploading k0s binary to host", slog.String("ID", s.id)))

	hk := HostGetK0s(h)
	hf := exec.HostGetFiles(h)
	hp := exec.HostGetPlatform(h)

	url := s.c.config.DownloadURL(hp.Arch(ctx))
	fs, fn, ferr := s.d.Download(ctx, url)
	if ferr != nil {
		return ferr
	}

	fnb := hk.getConfigurer().K0sBinaryPath()

	if _, fierr := hf.Stat(ctx, fnb, exec.ExecOptions{Sudo: true}); fierr != nil { // @TODO we have to compare local and uploaded file info
		slog.InfoContext(ctx, fmt.Sprintf("%s: uploading k0s binary: %s -> %s (cache: %s)", h.Id(), url, fnb, fn))
		if err := hf.Upload(ctx, fs, fnb, 0750, exec.ExecOptions{Sudo: true}); err != nil {
			return err
		}
	} else {
		slog.InfoContext(ctx, fmt.Sprintf("%s: not uploading k0s binary, as it already exists on the %s", h.Id(), fn))
	}

	return nil
}

func (s installK0sStep) configureK0s(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: writing config file", slog.String("ID", s.id)))
	hk := HostGetK0s(h)

	cfgNew, cfberr := s.configFor(ctx, h)
	if cfberr != nil {
		return fmt.Errorf("failed to build k0s config for %s: %w", h, cfberr)
	}

	hf := exec.HostGetFiles(h)

	cfgior := strings.NewReader(cfgNew)
	hcp := hk.K0sConfigPath()
	if err := hf.Upload(ctx, cfgior, hcp, 0600, exec.ExecOptions{Sudo: true}); err != nil {
		return err
	}

	return nil
}

func (s installK0sStep) initializeK0s(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: initializing k0s", slog.String("ID", s.id)))
	he := exec.HostGetExecutor(h)
	hk := HostGetK0s(h)

	if s.c.config.DynamicConfig || (hk.InstallFlags.Include("--enable-dynamic-config") && hk.InstallFlags.GetValue("--enable-dynamic-config") != "false") {
		s.c.config.DynamicConfig = true
		hk.InstallFlags.AddOrReplace("--enable-dynamic-config")
	}

	if s.Force {
		slog.WarnContext(ctx, fmt.Sprintf("%s: --force given, using k0s install with --force", h.Id()))
		hk.InstallFlags.AddOrReplace("--force=true")
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: installing k0s controller", s.Id()))
	cmd, err := hk.K0sInstallCommand(ctx)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, fmt.Sprintf("install first k0s controller using `%s`", strings.ReplaceAll(cmd, hk.getConfigurer().K0sBinaryPath(), "k0s")))
	if _, e, err := he.Exec(ctx, cmd, nil, exec.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: %s", err.Error(), e)
	}

	hk.Metadata.K0sInstalled = true

	return nil
}

func (s installK0sStep) startK0s(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: starting k0s", slog.String("ID", s.id)))

	he := exec.HostGetExecutor(h)
	hk := HostGetK0s(h)

	if err := he.ServiceEnable(ctx, []string{hk.K0sServiceName()}); err != nil {
		return err
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: waiting for the k0s service to start", h.Id()))
	if err := retry.Timeout(context.TODO(), retry.DefaultTimeout, func(ctx context.Context) error {
		return he.ServiceIsRunning(ctx, []string{hk.K0sServiceName()})
	}); err != nil {
		return err
	}

	//port := 6443
	//slog.InfoContext(ctx, fmt.Sprintf("%s: waiting for kubernetes api to respond", h))
	//if err := retry.Timeout(context.TODO(), retry.DefaultTimeout, node.KubeAPIReadyFunc(h, port)); err != nil {
	//	return err
	//}

	//if p.IsWet() && p.Config.Spec.K0s.DynamicConfig {
	//	if err := retry.Timeout(context.TODO(), retry.DefaultTimeout, node.K0sDynamicConfigReadyFunc(h)); err != nil {
	//		return fmt.Errorf("dynamic config reconciliation failed: %w", err)
	//	}
	//}

	hk.Metadata.K0sRunningVersion = &s.c.config.Version
	hk.Metadata.K0sBinaryVersion = &s.c.config.Version
	hk.Metadata.Ready = true

	//if id, err := p.Config.Spec.K0s.GetClusterID(h); err == nil {
	//		p.Config.Spec.K0s.Metadata.ClusterID = id
	//		p.SetProp("clusterID", id)
	//	}
	//}

	return nil
}

func (s installK0sStep) configFor(ctx context.Context, h *host.Host) (string, error) {
	hn, _ := network.HostGetNetwork(h).Network(ctx)

	cfg := s.c.config.K0sConfig

	addr := hn.PrivateAddress
	sans := []string{hn.PrivateAddress}

	cfg.DigMapping("spec", "api")["address"] = addr
	addUnlessExist(&sans, addr)

	oldsans := cfg.Dig("spec", "api", "sans")
	switch oldsans := oldsans.(type) {
	case []interface{}:
		for _, v := range oldsans {
			if s, ok := v.(string); ok {
				addUnlessExist(&sans, s)
			}
		}
	case []string:
		for _, v := range oldsans {
			addUnlessExist(&sans, v)
		}
	}

	cs, cserr := s.c.GetControllerHosts(ctx)
	if cserr != nil {
		return "", cserr
	}

	for _, c := range cs {
		cn, cnerr := network.HostGetNetwork(c).Network(ctx)
		if cnerr != nil {

		}
		addUnlessExist(&sans, cn.PrivateAddress)
	}
	addUnlessExist(&sans, "127.0.0.1")
	cfg.DigMapping("spec", "api")["sans"] = sans

	if cfg.Dig("spec", "storage", "etcd", "peerAddress") != nil || hn.PrivateAddress != "" {
		cfg.DigMapping("spec", "storage", "etcd")["peerAddress"] = addr
	}

	if _, ok := cfg["apiVersion"]; !ok {
		cfg["apiVersion"] = "k0s.k0sproject.io/v1beta1"
	}

	if _, ok := cfg["kind"]; !ok {
		cfg["kind"] = "ClusterConfig"
	}

	c, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: host configuration built", h.Id()), slog.Any("cfg", cfg), slog.String("cfgyaml", string(c)))

	return fmt.Sprintf("# generated-by-launchpad %s\n%s", time.Now().Format(time.RFC3339), c), nil
}

func addUnlessExist(slice *[]string, s string) {
	var found bool
	for _, v := range *slice {
		if v == s {
			found = true
			break
		}
	}
	if !found {
		*slice = append(*slice, s)
	}
}
