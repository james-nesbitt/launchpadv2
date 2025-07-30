package k0s_test

import (
	"testing"

	"github.com/Mirantis/launchpad/mirantis/product/k0s"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/mock"
	"gopkg.in/yaml.v3"
)

func Test_HostPluginDecode(t *testing.T) {
	ctx := t.Context()
	ys := `
role: controller
reset: false
datadir: /var/run/k0s/data
`
	h := host.NewHost("test")
	h.AddPlugin(mock.NewMockHostPlugin(h, nil))

	var yn yaml.Node

	if err := yaml.Unmarshal([]byte(ys), &yn); err != nil {
		t.Fatal("failed to unmarshall k0s plugin yaml")
	}

	pf := k0s.HostPluginFactory{}

	p, perr := pf.HostPluginDecode(ctx, h, yn.Decode)
	if perr != nil {
		t.Fatalf("failed to decode k0s host plugin: %s", perr.Error())
	}

	if p == nil {
		t.Error("plugin decode failed")
	}
}
