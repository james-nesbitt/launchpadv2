package rig_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/implementation/rig"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

func Test_HostPluginRoleTypeSanity(t *testing.T) {
	ctx := context.Background()

	pf := rig.HostPluginFactory{}

	h := host.NewHost("test")
	p := pf.HostPlugin(ctx, h)

	h.AddPlugin(p)

	if e := exec.HostGetExecutor(h); e == nil {
		t.Errorf("%s: K0S not an exec host-plugin", h.Id())
	}

	if f := exec.HostGetFiles(h); f == nil {
		t.Errorf("%s: K0S not a files host-plugin", h.Id())
	}

	if n := network.HostGetNetwork(h); n == nil {
		t.Errorf("%s: K0S not a files host-plugin", h.Id())
	}

	if p := exec.HostGetPlatform(h); p == nil {
		t.Errorf("%s: K0S not a platform host-plugin", h.Id())
	}
}
