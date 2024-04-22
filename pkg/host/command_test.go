package host_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/host"
)

func Test_CommandSanity(t *testing.T) {
	var _ action.CommandHandler = &host.HostsComponent{}
}
