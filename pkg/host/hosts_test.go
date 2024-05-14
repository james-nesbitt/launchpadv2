package host_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
)

func Test_HostsEachSanity(t *testing.T) {
	ctx := context.Background()

	hs := host.NewHosts(
		host.NewHost("one"),
		host.NewHost("two"),
		host.NewHost("three"),
		host.NewHost("four"),
		host.NewHost("five"),
		host.NewHost("six"),
	)

	var i int = 0
	err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		i = i + 1
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error occurred in Each exec: %s", err.Error())
	}

	if i != 6 {
		t.Error("wrong number of execs occurred in Hosts.Each")
	}
}

func Test_HostsEachError(t *testing.T) {
	ctx := context.Background()

	hs := host.NewHosts(
		host.NewHost("one"),
	)

	err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		return errors.New("an error")
	})
	if err == nil {
		t.Errorf("expected error occurred in Each not returned")
	}
}
