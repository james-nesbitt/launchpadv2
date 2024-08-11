//go:build integration

package integration_test

import (
	"context"
	"testing"

	// Register v2 spec handler
	"github.com/Mirantis/launchpad/pkg/project"
	_ "github.com/Mirantis/launchpad/pkg/config/v2_0"
	_ "github.com/Mirantis/launchpad/pkg/config/v2_1"
	// Register Host handlers
	_ "github.com/Mirantis/launchpad/pkg/implementation/rig"
	// Register actual product handlers for testing
	_ "github.com/Mirantis/launchpad/pkg/product/mcr"
	_ "github.com/Mirantis/launchpad/pkg/product/mke3"
	_ "github.com/Mirantis/launchpad/pkg/product/msr2"
	// Register mock stuff
	_ "github.com/Mirantis/launchpad/pkg/mock"
	
	"github.com/Mirantis/launchpad/integration"
	"github.com/Mirantis/launchpad/pkg/config"
)

func Test_LegacyDiscoverCommand_20(t *testing.T) {
	ctx := context.Background()
	y := integration.IntegrationTestYamlDoNotCommit

	cl, err := config.ConfigFromYamllBytes([]byte(y))
	if err != nil {
		t.Fatalf("Error occured unarshalling yaml: %s \nYAML:\b%s", err.Error(), y)
	}

	if valerr := cl.Validate(ctx); valerr != nil {
		t.Errorf("project validation error: %s", valerr.Error())
	}

	cmd, err := cl.Command(ctx, project.CommandKeyApply)
	if err != nil {
		t.Fatalf("project command build error: %s", err.Error())
	}

	if err := cmd.Validate(ctx); err != nil {
		t.Fatalf("project command [%s] validation failed: %s", cmd.Key, err.Error())
	}

	if err := cmd.Run(ctx); err != nil {
		t.Errorf("project command [%s] execution failed: %s", cmd.Key, err.Error())
	}
}

func Test_LegacyDiscoverCommand_20(t *testing.T) {
	ctx := context.Background()
	y := integration.IntegrationTestYamlDoNotCommit

	cl, err := config.ConfigFromYamllBytes([]byte(y))
	if err != nil {
		t.Fatalf("Error occured unarshalling yaml: %s \nYAML:\b%s", err.Error(), y)
	}

	if valerr := cl.Validate(ctx); valerr != nil {
		t.Errorf("project validation error: %s", valerr.Error())
	}

	cmd, err := cl.Command(ctx, project.CommandKeyApply)
	if err != nil {
		t.Fatalf("project command build error: %s", err.Error())
	}

	if err := cmd.Validate(ctx); err != nil {
		t.Fatalf("project command [%s] validation failed: %s", cmd.Key, err.Error())
	}

	if err := cmd.Run(ctx); err != nil {
		t.Errorf("project command [%s] execution failed: %s", cmd.Key, err.Error())
	}
}
