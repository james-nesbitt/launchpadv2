package v20_test

/**
 * Centralized testing for various configuration options
 */

import (
	"testing"

	// Register mock Host and Product handlers
	_ "github.com/Mirantis/launchpad/pkg/host/mock"
	_ "github.com/Mirantis/launchpad/pkg/product/mock"
	// Register actual product handlers for testing
	_ "github.com/Mirantis/launchpad/pkg/product/k0s"
	_ "github.com/Mirantis/launchpad/pkg/product/mcr"
	_ "github.com/Mirantis/launchpad/pkg/product/mke3"
	_ "github.com/Mirantis/launchpad/pkg/product/mke4"
	_ "github.com/Mirantis/launchpad/pkg/product/msr2"
	_ "github.com/Mirantis/launchpad/pkg/product/msr4"

	"github.com/Mirantis/launchpad/pkg/cluster"
	v2_0 "github.com/Mirantis/launchpad/pkg/config/v2_0"
	"github.com/Mirantis/launchpad/pkg/util/decode"
)

func Test_DecodeSanity(t *testing.T) {
	cl := cluster.Cluster{}
	sy := `
hosts:
  dummy:
    handler: mock
    roles:
    - dummy
products:
  dummy:
    handler: mock
  mock:
    dummy: nothing
`

	df := decode.DecodeTestYaml([]byte(sy))

	if err := v2_0.Decode(&cl, df); err != nil {
		t.Errorf("2.0 Spec decode failed unexpectadly: %s", err.Error())
	}

	if h, ok := cl.Hosts["dummy"]; !ok {
		t.Error("dummy host was not registered")
	} else if !h.HasRole("dummy") {
		t.Error("dummy host did not have the dummy role")
	}
}

func TestConfig_CurrentGen(t *testing.T) {
	cl := cluster.Cluster{}
	cy := `
hosts:
  dummy:
    handler: mock
    roles:
    - dummy
products:
  mcr:
    version: 23.0.7
  mke3:
    version: 3.7.2
  msr2:
    version: 2.9.10
`

	df := decode.DecodeTestYaml([]byte(cy))

	if err := v2_0.Decode(&cl, df); err != nil {
		t.Errorf("CurrentGen 2.0 Spec decode failed unexpectadly: %s", err.Error())
	}

	if len(cl.Components) != 3 {
		t.Errorf("Wrong number of components: %+v", cl.Components)
	}

	if mke, ok := cl.Components["mke3"]; !ok {
		t.Error("K0s component didn't decode properly")
	} else if mke.Name() != "MKE3" {
		t.Errorf("K0s product has wrong name: %s", mke.Name())
	}
}

func TestConfig_NextGen(t *testing.T) {
	cl := cluster.Cluster{}
	cy := `
hosts:
  c1:
    handler: mock
    roles:
    - controller
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
		t.Fatalf("NextGen 2.0 Spec decode failed unexpectadly: %s", err.Error())
	}

	if len(cl.Components) != 3 {
		t.Errorf("Wrong number of components: %+v", cl.Components)
	}

	if k0s, ok := cl.Components["k0s"]; !ok {
		t.Error("K0s component didn't decode properly")
	} else if k0s.Name() != "K0S" {
		t.Errorf("K0s product has wrong name: %s", k0s.Name())
	}
}
