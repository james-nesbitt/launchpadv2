package mock_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_MockHostPluginFactory(t *testing.T) {
	ctx := context.Background()
	pf := mock.MockHostPluginFactory{}
	h := host.NewHost("test-host")

	hp := pf.HostPlugin(ctx, h)

	assert.Implements(t, (*host.HostPlugin)(nil), hp, "mock host plugin does not implement host plugin inferface")
}

func Test_MockHostPluginFactoryDecode(t *testing.T) {
	ctx := context.Background()
	pf := mock.MockHostPluginFactory{}
	h := host.NewHost("test-host")

	y := `
network:
  private_address: 10.0.0.1
  public_address: mock.example.org
  private_interface: mock0
`

	ynd := yaml.Node{}

	if err := yaml.Unmarshal([]byte(y), &ynd); err != nil {
		t.Error("could not unmarshal")
	}

	hp, hperr := pf.HostPluginDecode(ctx, h, ynd.Decode)

	require.Nil(t, hperr, "host plugin factory decode returned unexpected error")
	assert.Implements(t, (*host.HostPlugin)(nil), hp, "mock host plugin does not implement host plugin inferface")
}

func Test_MockHostPluginFactoryDecodeHost(t *testing.T) {
	ctx := context.Background()

	y := `
mock:
  network:
    private_address: 10.0.0.1
    public_address: mock.example.org
    private_interface: mock0
`

	dns := map[string]yaml.Node{}

	if err := yaml.Unmarshal([]byte(y), &dns); err != nil {
		t.Fatal("could not unmarshal")
	}

	ds := map[string]func(interface{}) error{}
	for id, dn := range dns {
		ds[id] = dn.Decode
	}

	assert.Greater(t, len(ds), 0, "no plugins decoded from yaml")

	id := "test-host"
	h, herr := host.DecodeHost(ctx, id, ds)

	assert.Nil(t, herr, "mock host decode returned unexpected error")
	assert.Equal(t, id, h.Id(), "mock host for decode has the wrong ID")

	slog.Error("host debug", slog.Any("h", h))

	mh := mock.HostGetMock(h)

	require.NotNil(t, mh, "decoded mock host could not mutate to a mock host")
}
