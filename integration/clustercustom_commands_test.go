//go:build integration

package integration_test

import (
	"context"
	"testing"

	// Register v2 spec handler.
	_ "github.com/Mirantis/launchpad/pkg/config/v2_0"
	_ "github.com/Mirantis/launchpad/pkg/config/v2_1"

	// Register Host handlers.
	_ "github.com/Mirantis/launchpad/pkg/implementation/rig"

	// register implementation components
	_ "github.com/Mirantis/launchpad/pkg/implementation/kubernetes/component"
	
	// Register legacy product handlers.
	_ "github.com/Mirantis/launchpad/pkg/product/mcr"
	_ "github.com/Mirantis/launchpad/pkg/product/mke3"
	_ "github.com/Mirantis/launchpad/pkg/product/msr2"
	_ "github.com/Mirantis/launchpad/pkg/product/k0s"
	_ "github.com/Mirantis/launchpad/pkg/product/mke4"
	_ "github.com/Mirantis/launchpad/pkg/product/msr3"
	_ "github.com/Mirantis/launchpad/pkg/product/msr4"
	_ "github.com/Mirantis/launchpad/pkg/product/mkex"

	// Register mock stuff, in case it gets used.
	_ "github.com/Mirantis/launchpad/pkg/mock"

	// Include integration pieces
	"github.com/Mirantis/launchpad/integration"
	"github.com/Mirantis/launchpad/pkg/config"
)

func Test_Custom_DiscoverCommand(t *testing.T) {
	ctx := context.Background()
	y := integration.IntegrationTestYamlDoNotCommit

	cl, err := config.ConfigFromYamllBytes([]byte(y))
	if err != nil {
		t.Fatalf("Error occured unarshalling yaml: %s \nYAML:\b%s", err.Error(), y)
	}

	if valerr := cl.Validate(ctx); valerr != nil {
		t.Errorf("project validation error: %s", valerr.Error())
	}

	cmd, err := cl.Command(ctx, project.CommandKeyDiscover)
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

func Test_Custom_ApplyCommand(t *testing.T) {
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

func Test_Custom_ResetCommand(t *testing.T) {
	ctx := context.Background()
	y := integration.IntegrationTestYamlDoNotCommit

	cl, err := config.ConfigFromYamllBytes([]byte(y))
	if err != nil {
		t.Fatalf("Error occured unarshalling yaml: %s \nYAML:\b%s", err.Error(), y)
	}

	if valerr := cl.Validate(ctx); valerr != nil {
		t.Errorf("project validation error: %s", valerr.Error())
	}

	cmd, err := cl.Command(ctx, project.CommandKeyReset)
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
