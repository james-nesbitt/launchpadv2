package host_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/host/mock"

	"github.com/Mirantis/launchpad/pkg/host"
)

func Test_DecodeHost(t *testing.T) {
	host.RegisterDecoder("dummy", func(func(interface{}) error) (host.Host, error) {
		h := mock.NewHost(
			[]string{"test"},
			nil,
		)
		return host.Host(h), nil
	})

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
