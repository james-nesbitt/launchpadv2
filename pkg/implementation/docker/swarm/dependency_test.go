package swarm_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	"github.com/Mirantis/launchpad/pkg/implementation/docker/swarm"
)

func Test_DependencySanity(t *testing.T) {
	ctx := context.Background()
	v := docker.Version{
		Version: "23.0.9",
	}
	cfg := swarm.Config{}
	ds := swarm.NewSwarm(cfg)

	// 1. a requirer will generate a requirement

	var r dependency.Requirement = swarm.NewDockerSwarmRequirement(
		"test",
		"my test requirement",
		v,
	)

	if r.Describe() != "my test requirement" {
		t.Errorf("Test requirement gave wrong description: %+v", r)
	}

	// 2. a provider will check if the requirement is of the matching type

	dsr, ok := r.(swarm.DockerSwarmRequirement)
	if !ok {
		t.Fatalf("Couldn't convert requirement to our type: %+v", r)
	}

	// 3. now the dependency provider will get the needed docker version data

	v2 := dsr.NeedsDockerSwarm(ctx)

	if v2.Version != v.Version {
		t.Errorf("DockerSwarm requirement gave wrong version: %+v", v2)
	}

	// 4. with the version data, the provier would build its dependency as needed

	dsd := swarm.NewDockerSwarmDependency(
		"test",
		"my test dependency",
		func(ctx context.Context) (*swarm.Swarm, error) {
			// realistically, a provider would use the v (Version) and some other knowledge to build this
			return ds, nil
		},
	)

	var d dependency.Dependency = dsd

	// 5. provider will submit the dependency to the requirement

	if r.Matched(ctx) != nil {
		t.Errorf("DockerSwarm requirement says it is matched before we matched it: %+v", r)
	}

	if err := r.Match(d); err != nil {
		t.Errorf("DockerSwarm requirment returned an error when matched: %s", err.Error())
	}

	// 6. The source of the requirement will retrieve the dependency from the requiremnt,

	d2 := r.Matched(ctx)
	if d2 == nil {
		t.Fatalf("DockerSwarm requirement didn't return matched dependency: %+v", r)
	}

	if d2.Id() != "test" {
		t.Errorf("DockerSwarm Requirement produced the wrong dependency after being matched: %+v", d2)
	}

	// 7. Then the source can convert the dependency to the needed type, and retrieve the Docker client that it needs

	dsd2, ok := d2.(swarm.DockerSwarmDependency)
	if !ok {
		t.Errorf("DockerSwarm Requirement returned wrong kind of Dependency: %+v", dsd2)
	}

	ds2 := dsd2.ProvidesDockerSwarm(ctx)

	if ds2 != ds {
		t.Error("DockerHost Dependency returned wrong docker")
	}
}
