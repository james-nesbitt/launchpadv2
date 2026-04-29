package k0s_test

import (
	"testing"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/mirantis/product/k0s"
)

func Test_hostNeedsInstall(t *testing.T) {
	v130, _ := version.NewVersion("v1.30.0+k0s.0")
	v131, _ := version.NewVersion("v1.31.0+k0s.0")

	tests := []struct {
		name    string
		info    k0s.HostDiscovery
		desired version.Version
		want    bool
	}{
		{
			name:    "not running: needs install",
			info:    k0s.HostDiscovery{Running: false},
			desired: *v130,
			want:    true,
		},
		{
			name:    "running, nil version: needs install",
			info:    k0s.HostDiscovery{Running: true, RunningVersion: nil},
			desired: *v130,
			want:    true,
		},
		{
			name:    "running, same version: skip",
			info:    k0s.HostDiscovery{Running: true, RunningVersion: v130},
			desired: *v130,
			want:    false,
		},
		{
			name:    "running, different version: needs upgrade",
			info:    k0s.HostDiscovery{Running: true, RunningVersion: v130},
			desired: *v131,
			want:    true,
		},
		{
			name:    "running, newer version: needs downgrade",
			info:    k0s.HostDiscovery{Running: true, RunningVersion: v131},
			desired: *v130,
			want:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := k0s.HostNeedsInstall(tc.info, tc.desired)
			if got != tc.want {
				t.Errorf("HostNeedsInstall() = %v, want %v", got, tc.want)
			}
		})
	}
}
