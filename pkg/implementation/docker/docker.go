package dockerimplementation

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	dockertypesimage "github.com/docker/docker/api/types/image"
	dockertypesswarm "github.com/docker/docker/api/types/swarm"
	dockertypessystem "github.com/docker/docker/api/types/system"
)

type RunOptions struct {
	ShowOutput bool
	ShowError  bool
}

type DockerImplementation interface {
	//	@NOTE It would be nice to avoid this and use just the API commands to create run and delete a container.
	Run(ctx context.Context, args []string, ro RunOptions) (string, string, error)
	// Version retrieve the Docker Version from the remote server.
	Version(ctx context.Context) (map[string]dockertypes.Version, error)
	// Info retrieve the Docker VInfo from the remote server.
	Info(ctx context.Context) (dockertypessystem.Info, error)
	// ImagePull
	ImagePull(ctx context.Context, refStr string, options dockertypesimage.PullOptions) (io.ReadCloser, error)
	// NodeList retried the list of nodes in the swarm
	NodeList(ctx context.Context, options dockertypes.NodeListOptions) ([]dockertypesswarm.Node, error)
	// SwarmInit Initialize swarm
	SwarmInit(ctx context.Context, r dockertypesswarm.InitRequest) error
	// DockerInspect get info and join tokens
	//
	// @NOTE the docker cli does not have an equivalent so we have to build it from a couple
	//
	//	of cli calls.
	SwarmInspect(ctx context.Context) (dockertypesswarm.Swarm, error)
	// SwarmJoin join a swarm
	SwarmJoin(ctx context.Context, r dockertypesswarm.JoinRequest) error
	// SwarmLeave leave a swarm cluster
	SwarmLeave(ctx context.Context, force bool) error
}
