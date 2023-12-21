package config_test

/**
 * Centralized testing for various configuration options
 */

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/config"
)

func TestStandardConfig(t *testing.T) {
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
    mke:
      version: 4.7.2
    msr:
      version: 2.9.10
  cluster:
    prune: true
`

	if _, err := config.ConfigFromYamllBytes([]byte(cys)); err != nil {
		t.Errorf("Could not standard build: %s", err.Error())
	}
}
