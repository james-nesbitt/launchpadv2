package mcr

import (
	"context"
	"fmt"
	"log/slog"

	dockertypes "github.com/docker/docker/api/types"

	dockerhost "github.com/Mirantis/launchpad/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/host"
)

type swarmActivateStep struct {
	baseStep

	id string
}

func (s swarmActivateStep) Id() string {
	return fmt.Sprintf("%s:mcr-swarm-activate", s.id)
}

func (s *swarmActivateStep) Run(ctx context.Context) error {
	/**
	 * @NOTE a lot of this functionality should be moved to the host plugin
	 */

	slog.InfoContext(ctx, "running MCR swarm-activate step", slog.String("ID", s.Id()))

	mhs, mhsgerr := s.c.GetManagerHosts(ctx)
	if mhsgerr != nil {
		return fmt.Errorf("could not retrieve managers to join the swarm: %s", mhsgerr.Error())
	}

	// 1. Find a swarm project OR Create a new one

	var l *host.Host

	if dl, err := dockerhost.DiscoverLeader(ctx, mhs); err == nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: discovered as state leader", dl.Id()), slog.Any("leader", dl))
		l = dl
	} else if il, err := dockerhost.InitSwarm(ctx, mhs); err == nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: became new state leader from swarm init", il.Id()), slog.Any("leader", il))
		l = il
	} else {
		return fmt.Errorf("could not initialize swarm")
	}

	ld := dockerhost.HostGetDockerExec(l)

	li, lierr := ld.Info(ctx)
	if lierr != nil {
		return fmt.Errorf("%s: swarm join failed because leader docker info error: %s", l.Id(), lierr.Error())
	}
	si, sierr := ld.SwarmInspect(ctx)
	if sierr != nil {
		return fmt.Errorf("%s: swarm join failed because leader docker swarm inspect error: %s", l.Id(), sierr.Error())
	}
	ni, nierr := ld.NodeList(ctx, dockertypes.NodeListOptions{})
	if nierr != nil {
		return fmt.Errorf("%s: swarm join failed because leader docker swarm node list error: %s", l.Id(), nierr.Error())
	}

	// at this point the state should have swarm info populated
	// so all we need to do is to join the rest of the hosts
	if err := mhs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: swarm manager join", h.Id()), slog.Any("host", h))
		return dockerhost.JoinSwarm(ctx, h, li, si, ni, "manager")
	}); err != nil {
		return fmt.Errorf("error joining managers to the swarm: %s", err.Error())
	}

	whs, whsgerr := s.c.GetWorkerHosts(ctx)
	if whsgerr != nil {
		return fmt.Errorf("could not retrieve workers to join the swarm: %s", whsgerr.Error())
	}
	if err := whs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: swarm worker join", h.Id()), slog.Any("host", h))
		return dockerhost.JoinSwarm(ctx, h, li, si, ni, "worker")
	}); err != nil {
		return fmt.Errorf("error joining workers to the swarm: %s", err.Error())
	}

	return nil
}
