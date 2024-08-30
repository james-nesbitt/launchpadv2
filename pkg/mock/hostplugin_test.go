package mock_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HostGetMockNoMockPlugin(t *testing.T) {
	id := "test-host"
	h := host.NewHost(id)

	hm := mock.HostGetMock(h)

	assert.Nil(t, hm, "host indicates that it has a mock plugin, even though it does not")
}

func Test_MockHostPluginSanity(t *testing.T) {
	ctx := context.Background()
	p := mock.NewMockHostPlugin(nil, nil)

	// Network stuff
	assert.Implements(t, (*network.HostNetwork)(nil), p, "mock host plugin is not a network.HostNetwork plugin")
	_, nerr := p.Network(ctx)
	assert.Error(t, nerr, "mock plugin with nil network did not return an error")

	// non matched role check
	assert.False(t, p.RoleMatch("not-a-real-role"), "mock host plugin says it matches a role which it shouldn't")
}

func Test_MockHostPluginInAHost(t *testing.T) {
	id := "test-host"
	h := host.NewHost(id)
	p := mock.NewMockHostPlugin(h, &network.Network{})

	assert.Equal(t, mock.HostRoleMock, p.Id(), "mock host plugin returned the wrong ID")
	assert.Nil(t, p.Validate(), "mock host plugin validate failed")
	assert.True(t, p.RoleMatch(mock.HostRoleMock), "mock host plugin like to match the mock role")

	h.AddPlugin(p)
	hm := mock.HostGetMock(h)

	assert.NotNil(t, hm, "mock host plugin added to host did not allow the host to provide mock role functionality")
	assert.Equal(t, mock.HostRoleMock, hm.Id(), "mock host plugin from host gave wrong ID")
}

func Test_MockHostPluginNetworkNil(t *testing.T) {
	id := "test-host"
	h := host.NewHost(id)
	var n *network.Network = nil
	p := mock.NewMockHostPlugin(h, n)

	h.AddPlugin(p)
	hm := network.HostGetNetwork(h)

	require.Nil(t, hm, "mock host plugin without network added to host still thought it could provide network information")
}

func Test_MockHostPluginNetwork(t *testing.T) {
	ctx := context.Background()

	id := "test-host"
	h := host.NewHost(id)
	n := &network.Network{
		PrivateAddress:   "10.0.0.1",
		PublicAddress:    "mock.example.org",
		PrivateInterface: "mock0",
	}
	p := mock.NewMockHostPlugin(h, n)

	h.AddPlugin(p)
	hm := network.HostGetNetwork(h)

	require.NotNilf(t, hm, "mock host plugin added to host did not allow the host to provide network functionality")

	hn, hnerr := hm.Network(ctx)

	assert.Nil(t, hnerr, "mock host plugin network was not returned")
	assert.NotNil(t, hn, "mock host plugin network was nil")
	assert.Equal(t, *n, hn, "mock host plugin had the wrong network")
}
