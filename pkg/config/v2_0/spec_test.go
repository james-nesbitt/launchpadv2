package v20_test

/**
 * Centralized testing for various configuration options
 */

import (
	"context"
	"testing"

	// Register mock Host and Product handlers.
	_ "github.com/Mirantis/launchpad/pkg/mock"

	// Register actual product handlers for testing.
	_ "github.com/Mirantis/launchpad/pkg/product/k0s"
	_ "github.com/Mirantis/launchpad/pkg/product/mcr"
	_ "github.com/Mirantis/launchpad/pkg/product/mke3"
	_ "github.com/Mirantis/launchpad/pkg/product/mke4"
	_ "github.com/Mirantis/launchpad/pkg/product/msr2"
	_ "github.com/Mirantis/launchpad/pkg/product/msr4"

	v2_0 "github.com/Mirantis/launchpad/pkg/config/v2_0"
	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/Mirantis/launchpad/pkg/util/decode"
)

func Test_DecodeSanity(t *testing.T) {
	cl := project.Project{}
	sy := `
hosts:
  dummy:
    mock: {}
products:
  dummy:
    handler: mock
  mock:
    dummy: nothing
`

	df := decode.DecodeTestYaml([]byte(sy))

	if err := v2_0.Decode(&cl, df); err != nil {
		t.Errorf("2.0 Spec decode failed unexpectedly: %s", err.Error())
	}

	if err := cl.Validate(context.Background()); err != nil {
		t.Errorf("2.0 Spec decode project validate failed: %s", err.Error())
	}
}

func TestConfig_CurrentGen(t *testing.T) {
	cl := project.Project{}
	cy := `
hosts:
  dummy-manager:
    mcr:
      role: manager
  dummy-worker:
    mcr: {}
products:
  mcr:
    version: 23.0.10
  mke3:
    version: 3.7.2
  msr2:
    version: 2.9.10
`

	df := decode.DecodeTestYaml([]byte(cy))

	if err := v2_0.Decode(&cl, df); err != nil {
		t.Errorf("CurrentGen 2.0 Spec decode failed unexpectedly: %s", err.Error())
	}

	if err := cl.Validate(context.Background()); err != nil {
		t.Errorf("2.0 Spec decode project validate failed: %s", err.Error())
	}

	if len(cl.Components) != 4 { // 3 products and the hosts component
		t.Errorf("Wrong number of components: %+v", cl.Components)
	}

	// 	if mke, ok := cl.Components["mke3"]; !ok {
	// 		t.Error("MKE component didn't decode properly")
	// 	} else if mke.Name() != "mke3" {
	// 		t.Errorf("MKE product has wrong name: %s", mke.Name())
	// 	}

}

func TestConfig_NextGen(t *testing.T) {
	cl := project.Project{}
	cy := `
hosts:
  dummy-controller:
    k0s:
      role: controller
products:
  k0s:
    version: v1.28.4+k0s.0
  mke4:
    version: 4.0.0
  msr4:
    version: 4.0.0
`

	df := decode.DecodeTestYaml([]byte(cy))

	if err := v2_0.Decode(&cl, df); err != nil {
		t.Fatalf("NextGen 2.0 Spec decode failed unexpectedly: %s", err.Error())
	}

	if len(cl.Components) != 4 { // 3 products and the hosts component
		t.Errorf("Wrong number of components: %+v", cl.Components)
	}

	// if k0s, ok := cl.Components["k0s"]; !ok {
	// 	t.Error("K0s component didn't decode properly")
	// } else if k0s.Name() != "k0s" {
	// 	t.Errorf("K0s product has wrong name: %s", k0s.Name())
	// }
}
