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
