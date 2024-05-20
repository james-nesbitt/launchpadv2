package rig

import (
	"context"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host/exec"
)

// IsWindows.
func (h rigHostPlugin) IsWindows(ctx context.Context) bool {
	return h.rig.IsWindows()
}

func (h rigHostPlugin) Arch(ctx context.Context) string {
	arch, _, err := h.Exec(ctx, "uname -m", nil, exec.ExecOptions{})
	if err != nil {
		return ""
	}

	arch = strings.TrimSpace(arch)

	switch arch {
	case "x86_64":
		return "amd64"
	case "aarch64":
		return "arm64"
	case "armv7l", "armv8l", "aarch32", "arm32", "armhfp", "arm-32":
		return "arm"
	default:
		return arch
	}
}

// UserCacheDir returns the default root directory to use for user-specific cached data.
func (p *rigHostPlugin) UserCacheDir() string {
	fs := p.rig.FS()

	if cache := fs.Getenv("XDG_CACHE_HOME"); cache != "" {
		return cache
	}
	return fs.Join(fs.UserHomeDir(), ".cache")
}

// UserConfigDir returns the default root directory to use for user-specific configuration data.
func (p *rigHostPlugin) UserConfigDir() string {
	fs := p.rig.FS()

	if config := fs.Getenv("XDG_CONFIG_HOME"); config != "" {
		return config
	}
	return fs.Join(fs.UserHomeDir(), ".config")
}

// UserHomeDir returns the current user's home directory.
func (p *rigHostPlugin) UserHomeDir() string {
	fs := p.rig.FS()

	return fs.Getenv("HOME")
}
