package mock_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/mock"
)

func Test_RolesPlugin_PluginSanty(t *testing.T) {
	h := host.NewHost("test")
	var mp host.HostPlugin = mock.NewMockPlugin(h)

	h.AddPlugin(mp)

	hmp := h.MatchPlugin(mock.HostRoleMock)
	if hmp == nil {
		t.Error("host couldn't match mock plugin")
	}
}
