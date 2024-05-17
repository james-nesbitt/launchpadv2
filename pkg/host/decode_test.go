package host_test

import (
	"context"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/mock"
	_ "github.com/Mirantis/launchpad/pkg/host/mock"
)

func Test_HostV2DecodeMock(t *testing.T) {
	ctx := context.Background()
	id := "test"
	y := `
mock:
  network:
    private_address: 10.0.0.1
    public_address: mock.example.org
    private_interface: mock0
`

	hpsc := map[string]yaml.Node{} // container for the plugins
	if err := yaml.Unmarshal([]byte(y), &hpsc); err != nil {
		t.Fatalf("yaml unmarshalling fail: %s", err.Error())
	}

	hpds := map[string]func(interface{}) error{}
	for k, hpd := range hpsc {
		hpds[k] = hpd.Decode
	}

	h, err := host.DecodeHost(ctx, id, hpds)
	if err != nil {
		t.Fatalf("host decode failed: %s", err.Error())
	}

	if h == nil {
		t.Error("host decode gave a nil")
	}

	mp := h.MatchPlugin(mock.HostRoleMock)
	if mp == nil {
		t.Error("decoded host didn't have the mock plugin role satisfied")
	}
}
