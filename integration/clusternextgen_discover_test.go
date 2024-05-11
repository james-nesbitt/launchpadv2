//go:build integration
package integration_test

import (
	"context"
	"testing"

	// Register v2 spec handler
	"github.com/Mirantis/launchpad/pkg/action"
	_ "github.com/Mirantis/launchpad/pkg/config/v2_0"
	// Register Host handlers
	_ "github.com/Mirantis/launchpad/pkg/host/mock"
	_ "github.com/Mirantis/launchpad/pkg/host/rig"
	// Register actual product handlers for testing
	_ "github.com/Mirantis/launchpad/pkg/product/k0s"
	_ "github.com/Mirantis/launchpad/pkg/product/mke4"
	_ "github.com/Mirantis/launchpad/pkg/product/msr4"

	"github.com/Mirantis/launchpad/integration"
	"github.com/Mirantis/launchpad/pkg/config"
)

func Test_NextGenDiscoverCommand(t *testing.T) {
	ctx := context.Background()
	y := integration.IntegrationTestYamlNextGen

	cl, err := config.ConfigFromYamllBytes([]byte(y))
	if err != nil {
		t.Fatalf("Error occured unarshalling yaml: %s \nYAML:\b%s", err.Error(), y)
	}

	if valerr := cl.Validate(ctx); valerr != nil {
		t.Errorf("cluster validation error: %s", valerr.Error())
	}

	cmd, err := cl.Command(ctx, action.CommandKeyDiscover)
	if err != nil {
		t.Fatalf("cluster command build error: %s", err.Error())
	}

	if err := cmd.Validate(ctx); err != nil {
		t.Errorf("cluster command [%s] validation failed: %s", cmd.Key, err.Error())
	}

	if err := cmd.Run(ctx); err != nil {
		t.Errorf("cluster command [%s] execution failed: %s", cmd.Key, err.Error())
	}
}
