package k0s_test

import (
	"strings"
	"testing"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product/k0s"
)

func Test_k0sProductSanity(t *testing.T) {
	var k0sp component.Component = k0s.NewComponent("dummy", k0s.Config{
		Version: *version.MustParse("v1.28.4+k0s.0"),
	})

	if !strings.Contains(k0sp.Name(), "dummy") {
		t.Errorf("K0s component gave wrong name: %s", k0sp.Name())
	}
}
