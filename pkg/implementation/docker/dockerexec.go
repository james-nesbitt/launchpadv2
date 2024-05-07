package dockerimplementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
)

/**
 * Run Docker by executing the Docker cli program.
 *
 * @NOTE using the golang Docker SDK would be better
 *
 * By Design, what we do here is to mock the Docker
 * SDK so that we can switch over to it if we can ge
 * it working.
 * @see https://pkg.go.dev/github.com/docker/docker/client
 */

var (
	ErrDockerExecuteError = errors.New("error occurred running Docker command")
)

func NewDockerExec(executor func(ctx context.Context, cmd string, inr io.Reader) (string, string, error)) DockerExec {
	return DockerExec{
		executor: executor,
	}
}

type DockerExec struct {
	executor func(ctx context.Context, cmd string, inr io.Reader) (string, string, error)
}

func (de DockerExec) dockerCommand(ctx context.Context, args []string) (string, string, error) {
	cmd := strings.Join(append([]string{"docker"}, args...), " ")
	return de.executor(ctx, cmd, nil)
}

// ServerVersion retrieve the Docker Version from the remote server.
func (de DockerExec) ServerVersion(ctx context.Context) (types.Version, error) {
	var dv types.Version

	o, e, eerr := de.dockerCommand(ctx, []string{"version", "--format=json"})
	if eerr != nil {
		return dv, fmt.Errorf("%w; %s : %s", ErrDockerExecuteError, eerr, e)
	}

	if err := json.Unmarshal([]byte(o), &dv); err != nil {
		return dv, fmt.Errorf("%w; unmarshal error %s", ErrDockerExecuteError, err)
	}

	return dv, nil
}
