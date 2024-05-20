package k0s

import (
	"context"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

type configurer interface {
	// K0sBinaryPath returns the path to the k0s binary on the host
	K0sBinaryPath() string
	// K0sConfigPath returns the path to the k0s config file on the host
	K0sConfigPath() string
	// K0sJoinTokenPath returns the path to the k0s join token file on the host
	K0sJoinTokenPath() string
	// DataDirDefaultPath returns the path to the k0s data dir on the host
	DataDirDefaultPath() string
	// SetPath sets a path for a key
	SetPath(key, value string)
	// K0sCmdf can be used to construct k0s commands in sprintf style.
	K0sCmdf(template string, args ...interface{}) string

	K0sBinaryVersion(ctx context.Context, h *host.Host) (*version.Version, error)
	// K0sctlLockFilePath returns a path to a lock file
	K0sctlLockFilePath(ctx context.Context, h *host.Host) string
	// ReplaceK0sTokenPath replaces the config path in the service stub
	ReplaceK0sTokenPath(ctx context.Context, h *host.Host, spath string) error
	// FileContains returns true if a file contains the substring
	FileContains(ctx context.Context, h *host.Host, path, s string) bool
	// CommandExist returns true if the command exists
	CommandExist(ctx context.Context, h *host.Host, cmd string, opts exec.ExecOptions) bool
	// KubeconfigPath returns the path to a kubeconfig on the host
	KubeconfigPath(ctx context.Context, h *host.Host, dataDir string) string
	// KubectlCmdf returns a command line in sprintf manner for running kubectl on the host using the kubeconfig from KubeconfigPath
	KubectlCmdf(ctx context.Context, h *host.Host, dataDir, s string, args ...interface{}) string
	// HTTPStatus makes a HTTP GET request to the url and returns the status code or an error
	HTTPStatus(ctx context.Context, h *host.Host, url string) (int, error)
}
