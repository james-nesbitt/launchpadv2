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
	_ "github.com/Mirantis/launchpad/pkg/product/mcr"
	_ "github.com/Mirantis/launchpad/pkg/product/mke3"
	_ "github.com/Mirantis/launchpad/pkg/product/msr2"

	"github.com/Mirantis/launchpad/integration"
	"github.com/Mirantis/launchpad/pkg/config"
)

func Test_LegacyResetCommand(t *testing.T) {
	ctx := context.Background()
	y := integration.IntegrationTestYamlLegacy

	cl, err := config.ConfigFromYamllBytes([]byte(y))
	if err != nil {
		t.Fatalf("Error occured unarshalling yaml: %s \nYAML:\b%s", err.Error(), y)
	}

	if valerr := cl.Validate(ctx); valerr != nil {
		t.Errorf("cluster validation error: %s", valerr.Error())
	}

	cmd, err := cl.Command(ctx, action.CommandKeyReset)
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
