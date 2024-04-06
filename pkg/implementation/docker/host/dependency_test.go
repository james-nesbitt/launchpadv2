package dockerhost_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

func Test_DependencySanity(t *testing.T) {
	ctx := context.Background()
	v := docker.Version{
		Version: "23.0.9",
	}
	dh := dockerhost.NewDockerHosts(
		&host.Hosts{},
		v,
	)

	// 1. a requirer will generate a requirement

	var r dependency.Requirement = dockerhost.NewDockerHostsRequirement(
		"test",
		"my test requirement",
		v,
	)

	if r.Describe() != "my test requirement" {
		t.Errorf("Test requirement gave wrong description: %+v", r)
	}

	// 2. a provider will check if the requirement is of the matching type

	dhr2, ok := r.(dockerhost.DockerHostsRequirement)
	if !ok {
		t.Fatalf("Couldn't convert requirement to our type: %+v", r)
	}

	// 3. now the dependency provider will get the needed docker version data

	v2 := dhr2.NeedsDockerHost(ctx)

	if v2.Version != v.Version {
		t.Errorf("DockerHost requirement gave wrong version: %+v", v2)
	}

	// 4. with the version data, the provier would build its dependency as needed

	dhd := dockerhost.NewDockerHostsDependency(
		"test",
		"my test dependency",
		func(ctx context.Context) (*dockerhost.DockerHosts, error) {
			return dh, nil
		},
	)

	var dd dependency.Dependency = dhd

	// 5. provider will submit the dependency to the requirement

	if r.Matched(ctx) != nil {
		t.Errorf("DockerHost requirement says it is matched before we matched it: %+v", r)
	}

	if err := r.Match(dd); err != nil {
		t.Errorf("DockerHost requirment returned an error when matched: %s", err.Error())
	}

	// 6. The source of the requirement will retrieve the dependency from the requiremnt,

	dd2 := r.Matched(ctx)
	if dd2 == nil {
		t.Fatalf("Docker Host requirement didn't return matched dependency: %+v", r)
	}

	if dd2.Id() != "test" {
		t.Errorf("DockerHost Requirement produced the wrong dependency after being matched: %+v", dd2)
	}

	// 7. Then the source can convert the dependency to the needed type, and retrieve the Docker client that it needs

	dhd2, ok := dd2.(dockerhost.DockerHostsDependency)
	if !ok {
		t.Errorf("DockerHost Requirement returned wrong kind of Dependency: %+v", dhd2)
	}

	dh2 := dhd2.ProvidesDockerHost(ctx)

	if dh2 != dh {
		t.Error("DockerHost Dependency returned wrong docker")
	}
}
