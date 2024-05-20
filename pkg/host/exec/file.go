package exec

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"

	"github.com/Mirantis/launchpad/pkg/host"
)

// HostFiles
//
// Plugin can manage files and folders on a host
const (
	HostRoleFiles = "files"
)

// HostGetFiles get the host plugin that can execute.
//
// If the hosts doesn't have an appropriate plugin then
// nil is returned.
func HostGetFiles(h *host.Host) HostFiles {
	hgf := h.MatchPlugin(HostRoleFiles)
	if hgf == nil {
		slog.Info(fmt.Sprintf("%s: no matching plugin for files: %s", h.Id(), HostRoleFiles))
		return nil
	}

	hf, ok := hgf.(HostFiles)
	if !ok {
		slog.Info(fmt.Sprintf("%s: plugin is not a files plugin: %s", h.Id(), hgf.Id()))
		return nil
	}

	return hf
}

type HostFiles interface {
	// Stat return file info on a file, or an error
	Stat(ctx context.Context, path string, opts ExecOptions) (os.FileInfo, error)
	// FileExists confirm that a file exists
	FileExist(ctx context.Context, path string, opts ExecOptions) error
	// Upload content from a io.Reader to a path on the Host
	Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts ExecOptions) error
	// Rename a file
	Rename(ctx context.Context, old, new string, opts ExecOptions) error
}
