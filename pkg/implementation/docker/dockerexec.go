package dockerimplementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
func (de DockerExec) ServerVersion(ctx context.Context) (map[string]types.Version, error) {
	var dv map[string]types.Version

	o, e, eerr := de.dockerCommand(ctx, []string{"version", "--format=json"})
	if eerr != nil {
		return dv, fmt.Errorf("%w; %s : %s", ErrDockerExecuteError, eerr, e)
	}

	if len(o) == 0 {
		return dv, fmt.Errorf("%w: no version retrieved: `%s` / `%s`", ErrDockerExecuteError, o, e)
	}

	if err := json.Unmarshal([]byte(o), &dv); err != nil {
		return dv, fmt.Errorf("%w; unmarshal error %s `%s`", ErrDockerExecuteError, err, o)
	}

	if len(dv) == 0 {
		return dv, fmt.Errorf("%w: no version discovered: %+v", ErrDockerExecuteError, dv)
	}

	slog.ErrorContext(ctx, "DockerVersion", slog.Any("version", dv))
	return dv, nil
}

// SwarmInit Initialize swarm
func (de DockerExec) SwarmInit(ctx context.Context, installFlags []string) error {
	cmd := []string{"swarm", "init"}
	cmd = append(cmd, installFlags...)

	o, e, err := de.dockerCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("swarm init failed: %w :: %s", err.Error(), e)
	}

	slog.InfoContext(ctx, fmt.Sprintf("swarm init suceeded: %s", o))

	return nil
}

// SwarmJoinToken get a swarm join token
func (de DockerExec) SwarmJoinToken(ctx context.Context, joinType string, rotate bool) (string, error) {
	cmd := []string{"swarm", "join-token", "-q"}

	if rotate {
		cmd = append(cmd, "--rotate")
	}

	cmd = append(cmd, joinType)

	o, e, err := de.dockerCommand(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("swarm token fail: %s :: %s", err.Error(), e)
	}

	return o, nil
}

// SwarmJoin join a swarm
func (de DockerExec) SwarmJoin(ctx context.Context, leaderAddress string, token string, options []string) error {
	cmd := []string{"swarm", "join"}

	cmd = append(cmd, fmt.Sprintf("--token=%s", token))

	if len(options) > 0 {
		cmd = append(cmd, options...)
	}

	cmd = append(cmd, leaderAddress)

	return nil
}
