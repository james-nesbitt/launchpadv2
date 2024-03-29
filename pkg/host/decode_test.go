package host_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
	"gopkg.in/yaml.v3"
)

// --- Dummy host ---

// DummyHost a host type that does nothing, but can be used as a substitute when needed
type DummyHost struct {
	Roles       []string `yaml:"roles"`
	DecodeError error
}

func (dh DummyHost) Exec(ctx context.Context, cmd string, inr io.Reader) (string, string, error) {
	return "see err", "DummyHost was asked to exec", fmt.Errorf("DummyHost don't exec: %+v", dh)
}

func (dh DummyHost) HasRole(rs string) bool {
	for _, r := range dh.Roles {
		if r == rs {
			return true
		}
	}
	return false
}

func Test_DecodeHost(t *testing.T) {
	host.HostDecoders["dummy"] = func(d host.Decoder) (host.Host, error) {
		return host.Host(&DummyHost{
			Roles: []string{"test"},
		}), nil
	}

	dhd := func(dh interface{}) error {
		return nil
	}

	dh, err := host.DecodeHost("dummy", dhd)
	if err != nil {
		t.Errorf("DummyHost decode returned an unepected error: %s", err.Error())
	}

	if !dh.HasRole("test") {
		t.Errorf("host decode missing expected role: %+v", dh)
	}
}

type dummyHostContainer struct {
	t    *testing.T
	Host interface{} `yaml:"host"`

	yn chan *yaml.Node
}

func (dhc *dummyHostContainer) UnmarshalYAML(yn *yaml.Node) error {
	dhc.yn <- yn
	return nil
}

func Test_DecodeRig(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 1000000000) // don't block for more than a second
	dyn := make(chan *yaml.Node)
	hjs := `
host:
  roles:
  - test
`

	go func() {
		dhc := dummyHostContainer{yn: dyn, t: t}
		yaml.Unmarshal([]byte(hjs), &dhc)
	}()

	select {
	case <-ctx.Done(): // timeout if we messed anything up
		t.Error("yaml node channel didn't received signal. Timeout.")
		close(dyn) // don't leave the channel open

	case py := <-dyn:
		h, err := host.DecodeHost("dummy", py.Decode)
		if err != nil {
			t.Errorf("host decode unexpected error: %s", err.Error())
		} else if h == nil {
			t.Errorf("Decode gave a nil Host")
		} else if !h.HasRole("test") {
			t.Errorf("host decode missing expected role: %+v", h)
		}
	}
}
