package host

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/action/stepped"
)

func Test_ActionSanity(t *testing.T) {
	var _ stepped.Step = discoverStep{}
}
