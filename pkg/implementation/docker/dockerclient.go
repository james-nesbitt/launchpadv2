package dockerimplementation

import (
	dockerclient "github.com/docker/docker/client"
)

/**
 * @TODO Using the DockerClient is the better way to go, over the
 * the exec client, but for our MCR case, it requires abstracting
 * the net.Conn for the http transport,
 */

// NewDockerClient Constructor.
func NewDockerClient() *DockerClient {
	return &DockerClient{}
}

// DockerClient Docker implementation that uses the golang Docker client code.
type DockerClient struct {
	state State

	dockerclient.Client
}

// State for the docker implementation on hosts.
type State struct {
}
