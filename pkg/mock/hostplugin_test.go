package mock_test

import (
	"context"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_MockDecode(t *testing.T) {
	ctx := context.Background()
	y := `
mock:
  network:
    private_address: 10.0.0.1
    public_address: mock.example.org
    private_interface: mock0
`

	ynd := yaml.Node{}

	if err := yaml.Unmarshal([]byte(y), &ynd); err != nil {
		t.Error("could not unmarshal")
	}

	pf := mock.MockHostPluginFactory{}
	h := host.NewHost("test")

	hp, hpderr := pf.HostPluginDecode(ctx, h, ynd.Decode)
	if hpderr != nil {
		t.Fatal("failed to decode mock host plugin")
	}

	if hp.Id() != mock.HostRoleMock {
		t.Fatalf("mock has wrong id: %s", hp.Id())
	}

	if err := hp.Validate(); err != nil {
		t.Fatalf("mock validator produced an unexpected error: %s", err.Error())
	}

}

func Test_HostGetMock(t *testing.T) {
	h := host.NewHost("test")
	var mp host.HostPlugin = mock.NewMockHostPlugin(h)

	h.AddPlugin(mp)

	hgmp := mock.HostGetMock(h)
	if hgmp == nil {
		t.Errorf("couldn't get the mock out of the host: %+v", h)
	}
}

func Test_PluginSanty(t *testing.T) {
	h := host.NewHost("test")
	var mp host.HostPlugin = mock.NewMockHostPlugin(h)

	h.AddPlugin(mp)

	hmp := h.MatchPlugin(mock.HostRoleMock)
	if hmp == nil {
		t.Error("host couldn't match mock plugin")
	}
}
