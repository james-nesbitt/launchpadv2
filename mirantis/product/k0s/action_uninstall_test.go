package k0s_test

import (
	"testing"

	"github.com/Mirantis/launchpad/mirantis/product/k0s"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/mock"
)

// Test_hostReset_noK0sPlugin verifies that a host without the k0s plugin
// returns nil from HostGetK0s instead of panicking.
func Test_hostReset_noK0sPlugin(t *testing.T) {
	h := host.NewHost("test-host")
	h.AddPlugin(mock.NewMockHostPlugin(h, nil))

	if k0s.HostGetK0s(h) != nil {
		t.Fatal("expected nil k0s plugin for host without k0s configured")
	}
}

// Test_HostGetK0s_nilForUnknownHost verifies HostGetK0s returns nil for a
// host that has no plugins at all.
func Test_HostGetK0s_nilForUnknownHost(t *testing.T) {
	h := host.NewHost("bare-host")
	if k0s.HostGetK0s(h) != nil {
		t.Error("expected nil for host with no plugins")
	}
}
