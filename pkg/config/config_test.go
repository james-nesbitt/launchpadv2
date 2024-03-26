package config_test

/**
 * Centralized testing for various configuration options
 */

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/config"
)

func TestConfig_CurrentGen(t *testing.T) {
	cys := `
apiVersion: launchpad.mirantis.com/v2.0 
kind: cluster
metadata:
  name: test
spec:
  hosts:
  - role: manager
  products:
    mcr:
      version: 23.0.7
    mke3:
      version: 3.7.2
    msr2:
      version: 2.9.10
  cluster:
    prune: true
`

	if _, err := config.ConfigFromYamllBytes([]byte(cys)); err != nil {
		t.Errorf("Could not parse current gen: %s", err.Error())
	}
}

func TestConfig_NextGen(t *testing.T) {
	cys := `
apiVersion: launchpad.mirantis.com/v2.0 
kind: cluster
metadata:
  name: test
spec:
  hosts:
  - role: controller
  products:
    k0s:
      version: v1.28.4+k0s.0
    mke4:
      version: 4.0.0
    msr4:
      version: 4.0.0
  cluster:
    prune: true
`

	if _, err := config.ConfigFromYamllBytes([]byte(cys)); err != nil {
		t.Errorf("Could not parse next gen: %s", err.Error())
	}
}
