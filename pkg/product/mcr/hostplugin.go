package mcr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	dockertypessystem "github.com/docker/docker/api/types/system"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

/**
 * MCR has a HostPlugin which allows MCR specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleMCR = "mcr"
)

var (
	// ManagerRoles the Host roles accepted for managers.
	ManagerRoles = []string{"manager"}

	MCRInstallerPath = "mcr-install"
	MCRServices      = []string{"docker", "containerd"}
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMCR, &hostPluginFactory{})
}

type hostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin.
func (pf *hostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &hostPlugin{
		h: h,
	}
	pf.ps = append(pf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin.
func (pf *hostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &hostPlugin{h: h}

	if err := d(p); err != nil {
		return p, err
	}

	pf.ps = append(pf.ps, p)

	return p, nil
}

// Get the MCR plugin from a Host.
func HostGetMCR(h *host.Host) *hostPlugin {
	hgmcr := h.MatchPlugin(HostRoleMCR)
	if hgmcr == nil {
		return nil
	}

	hmcr, ok := hgmcr.(*hostPlugin)
	if !ok {
		return nil
	}

	return hmcr
}

// hostPlugin
//
// Implements no generic plugin interfaces.
type hostPlugin struct {
	h                *host.Host
	SwarmRole        string `yaml:"role"`
	ShouldSudoDocker bool   `yaml:"sudo_docker"`
	DaemonJson       string `yaml:"daemon_json"`
}

func (_ hostPlugin) Id() string {
	return "mcr"
}

func (_ hostPlugin) Validate() error {
	return nil
}

func (_ hostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMCR:
		return true
	case dockerhost.HostRoleDockerExec:
		// pkg/implementation/docker/host
		return true
	}

	return false
}

// DockerExecOptions meet the dockerhost.HostDockerExec interface.
func (hp hostPlugin) DockerExec() *dockerimplementation.DockerExec {
	e := exec.HostGetExecutor(hp.h)
	if e == nil {
		return nil
	}

	def := func(ctx context.Context, cmd string, i io.Reader, rops dockerimplementation.RunOptions) (string, string, error) {
		hopts := exec.ExecOptions{Sudo: hp.ShouldSudoDocker}

		if rops.ShowOutput {
			hopts.OutputLevel = "info"
		}
		if rops.ShowError {
			hopts.ErrorLevel = "warn"
		}

		return e.Exec(ctx, cmd, i, hopts)
	}

	return dockerimplementation.NewDockerExec(def)
}

func (hp hostPlugin) DockerConfig() string {
	return hp.DaemonJson
}

func (hp hostPlugin) IsManager() bool {
	for _, r := range ManagerRoles {
		if hp.SwarmRole == r {
			return true
		}
	}
	return false
}

func (hp hostPlugin) DockerInfo(ctx context.Context) (dockertypessystem.Info, error) {
	return hp.DockerExec().Info(ctx)
}

func (hp hostPlugin) DownloadMCRInstaller(ctx context.Context, c Config) error {
	ir := c.InstallURLLinux

	slog.InfoContext(ctx, fmt.Sprintf("%s: downloading MCR Installer: %s", hp.h.Id(), ir), slog.Any("host", hp.h))

	irs, igerr := http.Get(ir)
	if igerr != nil {
		return igerr
	}

	defer irs.Body.Close()

	if err := exec.HostGetFiles(hp.h).Upload(ctx, irs.Body, MCRInstallerPath, 0777, exec.ExecOptions{}); err != nil {
		return err
	}

	return nil
}

func (hp hostPlugin) RunMCRInstaller(ctx context.Context, c Config) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: running MCR Installer (%s)", hp.h.Id(), MCRInstallerPath), slog.Any("host", hp.h))

	cmd := fmt.Sprintf("DOCKER_URL=%s CHANNEL=%s VERSION=%s bash %s", c.RepoURL, c.Channel, c.Version, MCRInstallerPath)
	if _, e, err := exec.HostGetExecutor(hp.h).Exec(ctx, cmd, nil, exec.ExecOptions{Sudo: true}); err != nil {
		return fmt.Errorf("%s: mcr installer fail: %s \n %s", hp.h.Id(), err.Error(), e)
	}
	return nil
}

func (hp hostPlugin) UninstallMCR(ctx context.Context) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: uninstall MCR: %s", hp.h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", hp.h))

	if err := exec.HostGetExecutor(hp.h).RemovePackages(ctx, MCRUninstallPackages); err != nil {
		return err
	}

	return nil
}

func (hp hostPlugin) EnableMCRService(ctx context.Context) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: enabling MCR services: %s", hp.h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", hp.h))

	if err := exec.HostGetExecutor(hp.h).ServiceEnable(ctx, MCRServices); err != nil {
		return err
	}

	return nil
}

func (hp hostPlugin) RestartMCRService(ctx context.Context) error {
	if err := exec.HostGetExecutor(hp.h).ServiceRestart(ctx, []string{"docker"}); err != nil {
		return err
	}

	return nil
}

func (hp hostPlugin) DisableMCRService(ctx context.Context) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: disabling MCR services: %s", hp.h.Id(), strings.Join(MCRServices, ", ")), slog.Any("host", hp.h))

	if err := exec.HostGetExecutor(hp.h).ServiceDisable(ctx, MCRServices); err != nil {
		return err
	}

	return nil
}
