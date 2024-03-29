package host_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

// --- Testing Mocks ---

type TestingHost struct {
	exec_ioout, exec_ioerr string
	exec_err               error
	roles                  []string
}

func (th TestingHost) Exec(ctx context.Context, cmd string, inr io.Reader) (string, string, error) {
	return th.exec_ioout, th.exec_ioerr, th.exec_err
}

func (th TestingHost) HasRole(rs string) bool {
	for _, r := range th.roles {
		if r == rs {
			return true
		}
	}
	return false
}

func Test_MockHost(t *testing.T) {
	h1 := TestingHost{
		roles:      []string{"one", "two"},
		exec_ioout: "out",
		exec_ioerr: "err",
		exec_err:   nil,
	}

	if !h1.HasRole("two") {
		t.Errorf("Mock testinghost doesn't check roles properly: %+v", h1)
	}
	if !h1.HasRole("one") {
		t.Errorf("Mock testinghost doesn't check roles properly: %+v", h1)
	}

	io, ie, e := h1.Exec(context.Background(), "command", strings.NewReader("input"))
	if io != "out" || ie != "err" || e != nil {
		t.Errorf("Mock testinghost doesn't exec properly: %+v", h1)
	}
}

// --- Tests ---
func Test_HostProvides(t *testing.T) {
	ctx := context.Background()
	hs := host.Hosts{
		TestingHost{
			roles: []string{"test"},
		},
	}

	r1 := host.ClusterHostsRoleRequirementStandard{
		Description:  "r1 test role req",
		Role:         "test",
		MinimumCount: 1,
	}
	if err := hs.Provides(ctx, &r1); err != nil {
		t.Errorf("Hosts Provides returned an unexpected error: %s", err.Error())
	}
	if err := r1.Validate(ctx); err != nil {
		t.Errorf("requirement validation failed: %s", err.Error())
	} else if mhs, err := r1.RetrieveMatchingHosts(); err != nil {
		t.Errorf("Retrieve hosts failed: %s", err.Error())
	} else if len(mhs) == 0 {
		t.Errorf("Not enought hosts retrieved: %+v", r1)
	}

	r2 := host.ClusterHostsRoleRequirementStandard{
		Description:  "r2 test role req",
		Role:         "no-such-role",
		MinimumCount: 1,
	}
	if err := hs.Provides(ctx, &r2); err == nil {
		t.Error("Missing host role did not result in an error")
	} else if !errors.Is(err, dependency.ErrDependencyNotMet) {
		t.Errorf("Hosts Provides returned an unexpected error: %s", err.Error())
	}

}
