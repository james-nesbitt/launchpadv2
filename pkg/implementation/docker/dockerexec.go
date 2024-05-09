package dockerimplementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	dockertypesswarm "github.com/docker/docker/api/types/swarm"
	dockertypessystem "github.com/docker/docker/api/types/system"
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
	slog.DebugContext(ctx, "DOCKER COMMAND", slog.String("command", cmd))
	return de.executor(ctx, cmd, nil)
}

// Version retrieve the Docker Version from the remote server.
func (de DockerExec) Version(ctx context.Context) (map[string]dockertypes.Version, error) {
	var dv map[string]dockertypes.Version

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

	return dv, nil
}

// Info retrieve the Docker VInfo from the remote server.
func (de DockerExec) Info(ctx context.Context) (dockertypessystem.Info, error) {
	var di dockertypessystem.Info

	o, e, eerr := de.dockerCommand(ctx, []string{"info", "--format=json"})
	if eerr != nil {
		return di, fmt.Errorf("%w; %s : %s", ErrDockerExecuteError, eerr, e)
	}

	if len(o) == 0 {
		return di, fmt.Errorf("%w: no info retrieved: `%s` / `%s`", ErrDockerExecuteError, o, e)
	}

	if err := json.Unmarshal([]byte(o), &di); err != nil {
		return di, fmt.Errorf("%w; unmarshal error %s `%s`", ErrDockerExecuteError, err, o)
	}

	return di, nil
}

// SwarmInit Initialize swarm
func (de DockerExec) SwarmInit(ctx context.Context, r dockertypesswarm.InitRequest) error {
	cmd := []string{"swarm", "init"}

	if r.AdvertiseAddr != "" {
		cmd = append(cmd, fmt.Sprintf("--advertise-addr=%s", r.AdvertiseAddr))
	}
	if r.ListenAddr != "" {
		cmd = append(cmd, fmt.Sprintf("--listen-addr=%s", r.ListenAddr))
	}
	if r.DataPathAddr != "" {
		cmd = append(cmd, fmt.Sprintf("--data-path-addr=%s", r.DataPathAddr))
	}
	if r.Availability != "" {
		cmd = append(cmd, fmt.Sprintf("--availability=%s", r.DataPathAddr))
	}

	cmd = append(cmd, r.DefaultAddrPool...)

	_, e, err := de.dockerCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("swarm init failed: %w :: %s", err, e)
	}

	return nil
}

// DockerInspect get info and join tokens
//
// @NOTE the docker cli does not have an equivalent so we have to build it from a couple
//
//	of cli calls.
func (de DockerExec) SwarmInspect(ctx context.Context) (dockertypesswarm.Swarm, error) {
	s := dockertypesswarm.Swarm{}

	if i, err := de.Info(ctx); err != nil {
		return s, fmt.Errorf("manager swarm info fail: %s", err.Error())
	} else if i.Swarm.Cluster == nil {
		return s, fmt.Errorf("manager swarm info fail, no cluster info: %+v", i)
	} else {
		s.ClusterInfo = *(i.Swarm.Cluster)
	}

	tcmd := []string{"swarm", "join-token", "-q"}
	if o, e, err := de.dockerCommand(ctx, append(tcmd, "manager")); err != nil {
		return s, fmt.Errorf("manager swarm token fail: %s :: %s", err.Error(), e)
	} else {
		s.JoinTokens.Manager = strings.TrimSpace(o)
	}
	if o, e, err := de.dockerCommand(ctx, append(tcmd, "worker")); err != nil {
		return s, fmt.Errorf("worker swarm token fail: %s :: %s", err.Error(), e)
	} else {
		s.JoinTokens.Worker = strings.TrimSpace(o)
	}

	slog.DebugContext(ctx, "swarm inspect", slog.Any("swarm", s))

	return s, nil
}

// SwarmJoin join a swarm
func (de DockerExec) SwarmJoin(ctx context.Context, r dockertypesswarm.JoinRequest) error {
	cmd := []string{"swarm", "join"}

	cmd = append(cmd, fmt.Sprintf("--token=%s", r.JoinToken))

	if r.AdvertiseAddr != "" {
		cmd = append(cmd, fmt.Sprintf("--advertise-addrs=%s", r.AdvertiseAddr))
	}
	if r.DataPathAddr != "" {
		cmd = append(cmd, fmt.Sprintf("--data-path-addrs=%s", r.AdvertiseAddr))
	}
	if r.ListenAddr != "" {
		cmd = append(cmd, fmt.Sprintf("--listen-addrs=%s", r.ListenAddr))
	}
	if r.Availability != "" {
		cmd = append(cmd, fmt.Sprintf("--availability=%s", r.Availability))
	}

	cmd = append(cmd, r.RemoteAddrs...)

	_, e, err := de.dockerCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("swarm join fail: %s :: %s = %+v", err.Error(), e, r)
	}

	return nil
}

// SwarmLeave leave a swarm cluster
func (de DockerExec) SwarmLeave(ctx context.Context, force bool) error {
	cmd := []string{"swarm", "leave"}

	if force {
		cmd = append(cmd, "--force")
	}

	_, e, err := de.dockerCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("swarm leave fail: %s :: %s", err.Error(), e)
	}

	return nil
}
