package k0s

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/k0sproject/version"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
	"github.com/Mirantis/launchpad/pkg/util/download"
	"github.com/Mirantis/launchpad/pkg/util/flags"
)

/**
 * K0S has a HostPlugin which allows K0S specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleK0S = "k0s"
)

var (
	// K0SManagerHostRoles the Host roles accepted for managers.
	K0SManagerHostRoles = []string{"controller"}
)

// Get the K0S plugin from a Host.
func HostGetK0s(h *host.Host) *hostPlugin {
	hgk0s := h.MatchPlugin(HostRoleK0S)
	if hgk0s == nil {
		return nil
	}

	hk0s, ok := hgk0s.(*hostPlugin)
	if !ok {
		return nil
	}

	return hk0s
}

// HostConfig.
type HostConfig struct {
	Role             string            `yaml:"role"`
	Reset            bool              `yaml:"reset,omitempty"`
	ShouldDownload   bool              `yaml:"direct_download,omitempty"`
	DataDir          string            `yaml:"dataDir,omitempty"`
	Environment      map[string]string `yaml:"environment,flow,omitempty"`
	UploadBinary     bool              `yaml:"uploadBinary,omitempty"`
	K0sBinaryPath    string            `yaml:"k0sBinaryPath,omitempty"`
	K0sDataPath      string            `yaml:"k0sDataPath,omitempty"`
	K0sConfigPath    string            `yaml:"k0sConfigPath,omitempty"`
	K0sTokenPath     string            `yaml:"k0sTokenPath,omitempty"`
	K0sDownloadURL   string            `yaml:"k0sDownloadURL,omitempty"`
	InstallFlags     flags.Flags       `yaml:"installFlags,omitempty"`
	OSIDOverride     string            `yaml:"os,omitempty"`
	HostnameOverride string            `yaml:"hostname,omitempty"`
	NoTaints         bool              `yaml:"noTaints,omitempty"`
}

type hostState struct {
}

var (
	// global Download Queue.
	qd = download.NewQueueDownload(nil)
)

// hostPlugin.
type hostPlugin struct {
	h *host.Host
	c HostConfig
	s hostState
}

// Id uniquely identify the plugin.
func (p hostPlugin) Id() string {
	return fmt.Sprintf("%s:%s", p.h.Id(), "k0s")
}

// RoleMatch what host roles does this host plugin act.
func (p hostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleK0S:
		return true
	}

	return false
}

// Validate the plugin configuration.
func (p *hostPlugin) Validate() error {
	return nil
}

func (p hostPlugin) IsController() bool {
	return p.c.Role == "controller"
}

// Version retrieve plugin host k0s version (know error if k0s binary is missing).
func (p *hostPlugin) Version(ctx context.Context) (k0sversion, error) {
	cmd := p.k0sCommand([]string{"version", "--json"})

	eh := exec.HostGetExecutor(p.h)

	var v k0sversion

	o, e, err := eh.Exec(ctx, cmd, nil, exec.ExecOptions{Sudo: true})
	if err != nil {
		return v, k0sErrorAnalyze(e, fmt.Errorf("failed to retrieve version; %w", err))
	}

	if err := json.Unmarshal([]byte(o), &v); err != nil {
		return v, fmt.Errorf("failed to unmarshal version: %s", err.Error())
	}

	return v, nil
}

// Status retrieve plugin host status (known error if binary is missing, or if k0s is not running).
func (p *hostPlugin) Status(ctx context.Context) (k0sstatus, error) {
	cmd := p.k0sCommand([]string{"status", "--out=json"})

	eh := exec.HostGetExecutor(p.h)

	var s k0sstatus

	o, e, err := eh.Exec(ctx, cmd, nil, exec.ExecOptions{Sudo: true})
	if err != nil {
		return s, k0sErrorAnalyze(e, fmt.Errorf("failed to retrieve status; %w", err))
	}

	if err := json.Unmarshal([]byte(o), &s); err != nil {
		return s, fmt.Errorf("failed to unmarshal status: %s", err.Error())
	}

	return s, nil
}

// DownloadK0sBinary get k0s onto the plugin host by directly downloading it from the host.
func (p *hostPlugin) DownloadK0sBinary(ctx context.Context, v version.Version) error {
	hf := exec.HostGetFiles(p.h)
	hp := exec.HostGetPlatform(p.h)
	u := DownloadK0sURL(v, hp.Arch(ctx))
	d := p.k0sBinaryPath()

	if err := hf.Download(ctx, u, d, 0550, exec.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: host couldn't download the K0s binary: %s -> %s :: %s", p.h.Id(), u, d, err.Error())
	}
	return nil
}

// UploadK0sBinary get the k0s binary onto the plugin host by downloading it to the launchpad machine and then streaming it to the plugin host.
func (p *hostPlugin) UploadK0sBinary(ctx context.Context, v version.Version) error {
	hf := exec.HostGetFiles(p.h)
	hp := exec.HostGetPlatform(p.h)
	u := DownloadK0sURL(v, hp.Arch(ctx))
	d := p.k0sBinaryPath()

	rc, _, derr := qd.Download(ctx, u) // queue downloads so that we avoid repeating any downloads
	if derr != nil {
		return fmt.Errorf("could not download k0s binary, before uploading to host %s : %s -> %s : %s", p.h.Id(), u, d, derr.Error())
	}

	if err := hf.Upload(ctx, rc, d, 0750, exec.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: failed to upload k0s binary to host : %s", p.h.Id(), err.Error())
	}
	return nil
}

// BuildAndWriteK0sConfig Build host config and write it to the host as yaml.
func (p *hostPlugin) BuildAndWriteK0sConfig(ctx context.Context, basecfg K0sConfig, sans []string) error {
	c, err := p.BuildHostConfig(ctx, basecfg, sans)
	if err != nil {
		return err
	}
	return p.WriteK0sConfig(ctx, c)
}

// WriteK0sConfig to a plugin host as yaml.
func (p *hostPlugin) WriteK0sConfig(ctx context.Context, cfg K0sConfig) error {
	cfgbs, merr := yaml.Marshal(cfg)
	if merr != nil {
		return merr
	}

	cfgsr := strings.NewReader(string(cfgbs))

	hf := exec.HostGetFiles(p.h)

	if err := hf.Upload(ctx, cfgsr, p.k0sConfigPath(), 0440, exec.ExecOptions{Sudo: true}); err != nil {
		return err
	}
	return nil
}

var token_mu sync.Mutex

// GenerateToken generate a new k0s join token for a role
//
// @NOTE a mutex is used to ensure that only one token is generated at a time.
func (p *hostPlugin) GenerateToken(ctx context.Context, role string, expiry time.Duration) (string, error) {
	cmd := []string{
		"token",
		"create",
		fmt.Sprintf("--role=%s", role),
		fmt.Sprintf("--expiry=%s", expiry),
		fmt.Sprintf("--data-dir=%s", p.k0sDataPath()),
	}

	eh := exec.HostGetExecutor(p.h)

	token_mu.Lock()
	defer token_mu.Unlock()
	o, e, err := eh.Exec(ctx, p.k0sCommand(cmd), nil, exec.ExecOptions{Sudo: true})
	if err != nil {
		return o, fmt.Errorf("%s: error generating join token: %s :: %s", p.h.Id(), err.Error(), e)
	}
	return strings.TrimSpace(o), nil
}

// ActivateNewCluster Initialize the plugin host as the first controller in a new cluster
//
// @TODO do more stuff to check if the cluster is already running?
func (p *hostPlugin) InstallNewCluster(ctx context.Context, c Config) error {
	eh := exec.HostGetExecutor(p.h)

	s, serr := p.Status(ctx)
	if serr == nil {
		if s.Role != RoleController {
			return fmt.Errorf("%s: host cannot be treated as a new cluster leader, as it is a %s", p.h.Id(), s.Role)
		}
		return nil
	}

	for k, ss := range map[string][]string{
		"worker":     {ServiceWorker},
		"controller": {ServiceController},
	} {
		if eh.ServiceIsRunning(ctx, ss) == nil {
			if err := eh.ServiceDisable(ctx, ss); err != nil {
				return fmt.Errorf("%s: failed to stop K0s services: %s", p.h.Id(), k)
			}
		}
	}

	args := []string{
		"install",
		RoleController,
		"--debug",
		"--force",
		fmt.Sprintf("--data-dir=%s", p.k0sDataPath()),
		fmt.Sprintf("--config=%s", p.k0sConfigPath()),
	}

	if c.DynamicConfig {
		args = append(args, "--enable-dynamic-config")
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: installing leader: %s", p.h.Id(), args))

	_, e, ierr := eh.Exec(ctx, p.k0sCommand(args), nil, exec.ExecOptions{Sudo: true})
	if ierr != nil {
		return fmt.Errorf("%s: failed to initialize new cluster: %s :: %s", p.h.Id(), ierr.Error(), e)
	}
	slog.InfoContext(ctx, fmt.Sprintf("%s: installed as new leader", p.h.Id()))

	slog.InfoContext(ctx, fmt.Sprintf("%s: starting leader service: %s", p.h.Id(), ServiceController))
	if err := eh.ServiceEnable(ctx, []string{ServiceController}); err != nil {
		return fmt.Errorf("%s: failed to start leader k0s service: %s", p.h.Id(), err.Error())
	}

	return nil
}

// JoinCluster Join the plugin host to an existing k0s cluster.
func (p *hostPlugin) JoinCluster(ctx context.Context, l *host.Host, role string, c Config) error {
	lkh := HostGetK0s(l)
	eh := exec.HostGetExecutor(p.h)
	fh := exec.HostGetFiles(p.h)

	jt, terr := lkh.GenerateToken(ctx, role, JoinTokenExpireDuration)
	if terr != nil {
		return fmt.Errorf("%s: failed to get token from leader: %s", p.h.Id(), terr.Error())
	}

	if err := fh.Upload(ctx, strings.NewReader(jt), p.k0sTokenPath(), 0640, exec.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: failed to write join token to host: %s", p.h.Id(), err.Error())
	}

	args := []string{
		"install",
		role,
		"--debug",
		"--force",
		fmt.Sprintf("--token-file=%s", p.k0sTokenPath()),
		fmt.Sprintf("--data-dir=%s", p.k0sDataPath()),
	}

	if c.DynamicConfig {
		args = append(args, "--enable-dynamic-config")
	}

	// disable any running services
	for k, ss := range map[string][]string{
		"worker":     {ServiceWorker},
		"controller": {ServiceController},
	} {
		if eh.ServiceIsRunning(ctx, ss) == nil {
			if err := eh.ServiceDisable(ctx, ss); err != nil {
				return fmt.Errorf("%s: failed to stop K0s services: %s", p.h.Id(), k)
			}
		}
	}

	services := []string{}
	switch role {
	case RoleController:
		services = append(services, ServiceController)
		args = append(args, fmt.Sprintf("--config=%s", p.k0sConfigPath()))
	case RoleWorker:
		services = append(services, ServiceWorker)
	}

	_, e, err := eh.Exec(ctx, p.k0sCommand(args), nil, exec.ExecOptions{Sudo: true})
	if err != nil {
		return fmt.Errorf("%s: failed to join cluster: %s :: %s", p.h.Id(), err.Error(), e)
	}

	if err := eh.ServiceEnable(ctx, services); err != nil {
		return fmt.Errorf("%s: failed to enable&start K0s services", p.h.Id())
	}

	return nil
}

// K0sKubeconfigAdmin retrieve the admin kubeconfig from a k0s controller
//
//	@NOTE only works on a controller
func (p *hostPlugin) K0sKubeconfigAdmin(ctx context.Context) (string, error) {
	if p.c.Role != RoleController {
		return "", fmt.Errorf("%s: Can't retrieve kubeconfig from non controller node", p.h.Id())
	}

	args := []string{
		"kubeconfig",
		"admin",
	}

	eh := exec.HostGetExecutor(p.h)

	o, e, err := eh.Exec(ctx, p.k0sCommand(args), nil, exec.ExecOptions{Sudo: true})
	if err != nil {
		return "", fmt.Errorf("%s: error retrieving kubeconfig: %s", p.h.Id(), e)
	}

	return o, nil
}

// ---- Local helpers ----

func (p *hostPlugin) k0sCommand(args []string) string {
	return strings.Join(append([]string{p.k0sBinaryPath()}, args...), " ")
}

func (p *hostPlugin) k0sBinaryPath() string {
	if p.c.K0sBinaryPath != "" {
		return p.c.K0sBinaryPath
	}
	return DefaultK0sBinaryPath
}

func (p *hostPlugin) k0sConfigPath() string {
	if p.c.K0sConfigPath != "" {
		return p.c.K0sConfigPath
	}
	return DefaultK0sConfigPath
}

func (p *hostPlugin) k0sDataPath() string {
	if p.c.K0sDataPath != "" {
		return p.c.K0sDataPath
	}
	return DefaultK0sDataPath
}

func (p *hostPlugin) k0sTokenPath() string {
	if p.c.K0sTokenPath != "" {
		return p.c.K0sTokenPath
	}
	return DefaultK0sTokenPath
}

// BuildHostConfig Modify a passed base K0s config with host specific values, and including passed additional sans.
func (p *hostPlugin) BuildHostConfig(ctx context.Context, basecfg K0sConfig, sans []string) (K0sConfig, error) {
	var hcfg K0sConfig = basecfg
	slog.DebugContext(ctx, "base config", slog.Any("config", hcfg))

	addUnlessExist := func(slice *[]string, s string) {
		for _, v := range *slice {
			if v == s {
				return
			}
		}
		*slice = append(*slice, s)
	}

	hn := network.HostGetNetwork(p.h)
	n, nerr := hn.Network(ctx)
	if nerr != nil {
		return hcfg, nerr
	}

	var addr string
	if n.PrivateAddress != "" {
		addr = n.PrivateAddress
	} else {
		addr = n.PublicAddress
	}

	hcfg.Spec.API.Address = addr
	hcfg.Spec.Storage.Etcd.PeerAddress = addr
	addUnlessExist(&sans, addr)

	for _, s := range sans {
		addUnlessExist(&hcfg.Spec.API.Sans, s)
	}

	addUnlessExist(&hcfg.Spec.API.Sans, "127.0.0.1")

	if hcfg.Spec.API.K0sApiPort == 0 {
		hcfg.Spec.API.K0sApiPort = 9443
	}
	if hcfg.Spec.API.Port == 0 {
		hcfg.Spec.API.Port = 6443
	}
	//	if hcfg.Spec.Konnectivity.AdminPort == 0 {
	//		hcfg.Spec.Konnectivity.AdminPort = 8443
	//	}
	//	if hcfg.Spec.Konnectivity.AgentPort == 0 {
	//		hcfg.Spec.Konnectivity.AgentPort = 8443
	//	}

	return hcfg, nil
}
