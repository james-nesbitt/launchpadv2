package k0s

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"sync"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

// linuxConfigurer is a base module for various linux OS support packages
type linuxConfigurer struct {
	paths  map[string]string
	pathMu sync.Mutex
}

// NOTE The Linux struct does not embed rig/os.Linux because it will confuse
// go as the distro-configurers' parents embed it too. This means you can't
// add functions to base Linux package that call functions in the rig/os.Linux package,
// you can however write those functions in the distro-configurers.
// An example of this problem is the ReplaceK0sTokenPath function, which would like to
// call `l.ServiceScriptPath("kos")`, which was worked around here by getting the
// path as a parameter.

func (l *linuxConfigurer) initPaths() {
	if l.paths != nil {
		return
	}
	l.paths = map[string]string{
		"K0sBinaryPath": "/usr/local/bin/k0s",
		//"K0sBinaryPath": "~/k0s-v1.28.4+k0s.0-amd64",
		"K0sConfigPath":      "/etc/k0s/k0s.yaml",
		"K0sJoinTokenPath":   "/etc/k0s/k0stoken",
		"DataDirDefaultPath": "/var/lib/k0s",
	}
}

// K0sBinaryPath returns the path to the k0s binary on the host
func (l *linuxConfigurer) K0sBinaryPath() string {
	l.pathMu.Lock()
	defer l.pathMu.Unlock()

	l.initPaths()
	return l.paths["K0sBinaryPath"]
}

// K0sConfigPath returns the path to the k0s config file on the host
func (l *linuxConfigurer) K0sConfigPath() string {
	l.pathMu.Lock()
	defer l.pathMu.Unlock()

	l.initPaths()
	return l.paths["K0sConfigPath"]
}

// K0sJoinTokenPath returns the path to the k0s join token file on the host
func (l *linuxConfigurer) K0sJoinTokenPath() string {
	l.pathMu.Lock()
	defer l.pathMu.Unlock()

	l.initPaths()
	return l.paths["K0sJoinTokenPath"]
}

// DataDirDefaultPath returns the path to the k0s data dir on the host
func (l *linuxConfigurer) DataDirDefaultPath() string {
	l.pathMu.Lock()
	defer l.pathMu.Unlock()

	l.initPaths()
	return l.paths["DataDirDefaultPath"]
}

// SetPath sets a path for a key
func (l *linuxConfigurer) SetPath(key, value string) {
	l.pathMu.Lock()
	defer l.pathMu.Unlock()

	l.initPaths()
	l.paths[key] = value
}

// K0sCmdf can be used to construct k0s commands in sprintf style.
func (l *linuxConfigurer) K0sCmdf(template string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", l.K0sBinaryPath(), fmt.Sprintf(template, args...))
}

func (l *linuxConfigurer) K0sBinaryVersion(ctx context.Context, h *host.Host) (*version.Version, error) {
	o, _, err := exec.HostGetExecutor(h).Exec(ctx, l.K0sCmdf("version"), nil, exec.ExecOptions{})
	if err != nil {
		return nil, err
	}

	version, err := version.NewVersion(o)
	if err != nil {
		return nil, err
	}

	return version, nil
}

// K0sctlLockFilePath returns a path to a lock file
func (l *linuxConfigurer) K0sctlLockFilePath(ctx context.Context, h *host.Host) string {
	if ferr := exec.HostGetFiles(h).FileExist(ctx, "/run/lock", exec.ExecOptions{Sudo: true}); ferr == nil {
		return "/run/lock/k0sctl"
	}

	return "/tmp/k0sctl.lock"
}

var trailingNumberRegex = regexp.MustCompile(`(\d+)$`)

func trailingNumber(s string) (int, bool) {
	match := trailingNumberRegex.FindStringSubmatch(s)
	if len(match) > 0 {
		i, err := strconv.Atoi(match[1])
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

// ReplaceK0sTokenPath replaces the config path in the service stub
func (l *linuxConfigurer) ReplaceK0sTokenPath(ctx context.Context, h *host.Host, spath string) error {
	_, _, err := exec.HostGetExecutor(h).Exec(ctx, fmt.Sprintf("sed -i 's^REPLACEME^%s^g' %s", l.K0sJoinTokenPath(), spath), nil, exec.ExecOptions{})
	return err
}

// FileContains returns true if a file contains the substring
func (l *linuxConfigurer) FileContains(ctx context.Context, h *host.Host, path, s string) bool {
	_, _, err := exec.HostGetExecutor(h).Exec(ctx, fmt.Sprintf(`grep -q "%s" "%s"`, s, path), nil, exec.ExecOptions{Sudo: true})
	return err == nil
}

// KubeconfigPath returns the path to a kubeconfig on the host
func (l *linuxConfigurer) KubeconfigPath(ctx context.Context, h *host.Host, dataDir string) string {
	// if admin.conf exists, use that
	adminConfPath := path.Join(dataDir, "pki/admin.conf")
	if exec.HostGetFiles(h).FileExist(ctx, adminConfPath, exec.ExecOptions{Sudo: true}) == nil {
		return adminConfPath
	}
	return path.Join(dataDir, "kubelet.conf")
}

// KubectlCmdf returns a command line in sprintf manner for running kubectl on the host using the kubeconfig from KubeconfigPath
func (l *linuxConfigurer) KubectlCmdf(ctx context.Context, h *host.Host, dataDir, s string, args ...interface{}) string {
	return fmt.Sprintf(`env "KUBECONFIG=%s" %s`, l.KubeconfigPath(ctx, h, dataDir), l.K0sCmdf(`kubectl %s`, fmt.Sprintf(s, args...)))
}

// HTTPStatus makes a HTTP GET request to the url and returns the status code or an error
func (l *linuxConfigurer) HTTPStatus(ctx context.Context, h *host.Host, url string) (int, error) {
	o, _, err := exec.HostGetExecutor(h).Exec(ctx, fmt.Sprintf(`curl -kso /dev/null --connect-timeout 20 -w "%%{http_code}" "%s"`, url), nil, exec.ExecOptions{})
	if err != nil {
		return -1, err
	}
	status, err := strconv.Atoi(o)
	if err != nil {
		return -1, fmt.Errorf("invalid response: %s", err.Error())
	}

	return status, nil
}

// CommandExist returns true if the command exists
func (l *linuxConfigurer) CommandExist(ctx context.Context, h *host.Host, cmd string, opts exec.ExecOptions) bool {
	_, _, err := exec.HostGetExecutor(h).Exec(ctx, fmt.Sprintf(`command -v -- "%s" 2> /dev/null`, cmd), nil, opts)
	return err == nil
}
