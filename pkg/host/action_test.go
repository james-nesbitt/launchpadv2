package host

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
)

func Test_ActionSanity(t *testing.T) {
	var _ action.Step = discoverStep{}
}
