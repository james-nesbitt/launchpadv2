package mcr

import (
	"context"
	"fmt"
	"log/slog"

	dockertypesswarm "github.com/docker/docker/api/types/swarm"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

type swarmActivateStep struct {
	baseStep

	id string
}

func (s swarmActivateStep) Id() string {
	return fmt.Sprintf("%s:mcr-swarm-activate", s.id)
}

func (s *swarmActivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR swarm-activate step", slog.String("ID", s.Id()))

	mhs, mhsgerr := s.c.GetManagerHosts(ctx)
	if mhsgerr != nil {
		return fmt.Errorf("could not retrieve managers to join the swarm: %s", mhsgerr.Error())
	}

	// 1. Find a swarm cluster OR Create a new one

	var l *dockerhost.Host
	var swaas []string

	if dl, err := discoverLeader(ctx, mhs); err == nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: discovered as state leader", dl.Id()), slog.Any("leader", dl))
		l = dl
	} else if il, err := initSwarm(ctx, mhs); err == nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: became new state leader from swarm init", il.Id()), slog.Any("leader", il))
		l = il
	} else {
		return fmt.Errorf("could not initialize swarm")
	}

	sw, err := l.Docker(ctx).SwarmInspect(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve swarm info from state leader: %s", err.Error())
	}

	ln, lnerr := l.Network(ctx)
	if lnerr != nil {
		return fmt.Errorf("%s: could not retrieve leader network for swarm initialization: %s", l.Id(), lnerr.Error())
	}
	if ln.PrivateInterface == "" {
		return fmt.Errorf("%s: leader network did not contain an address we can use for swarm: %+v", l.Id(), ln)
	}
	slog.InfoContext(ctx, "setting swarm listen address", slog.Any("network", ln))
	swaas = append(swaas, ln.PrivateAddress)

	// at this point the state should have swarm info populated
	// so all we need to do is to join the rest of the hosts
	if err := mhs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		return joinSwarm(ctx, h, sw, swaas, "manager")
	}); err != nil {
		return fmt.Errorf("error joining managers to the swarm: %s", err.Error())
	}

	whs, whsgerr := s.c.GetWorkerHosts(ctx)
	if whsgerr != nil {
		return fmt.Errorf("could not retrieve workers to join the swarm: %s", whsgerr.Error())
	}
	if err := whs.Each(ctx, func(ctx context.Context, h *dockerhost.Host) error {
		return joinSwarm(ctx, h, sw, swaas, "worker")
	}); err != nil {
		return fmt.Errorf("error joining workers to the swarm: %s", err.Error())
	}

	return nil
}

func discoverLeader(ctx context.Context, mhs dockerhost.Hosts) (*dockerhost.Host, error) {
	for _, h := range mhs {
		i, ierr := h.Docker(ctx).Info(ctx)
		if ierr != nil {
			continue
		}

		if i.Swarm.ControlAvailable {
			slog.DebugContext(ctx, fmt.Sprintf("%s: swarm already active, this host can act as a leader.", h.Id()), slog.Any("host", h))
			return h, nil
		}
	}
	return nil, fmt.Errorf("No swarm leader discovered")
}

func initSwarm(ctx context.Context, mhs dockerhost.Hosts) (*dockerhost.Host, error) {
	for _, h := range mhs {
		n, nerr := h.Network(ctx)
		if nerr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: could not retrieve network for swarm initialization: %s", h.Id(), nerr.Error()))
			continue
		}

		ir := dockertypesswarm.InitRequest{
			AdvertiseAddr: n.PrivateAddress,
		}

		if err := h.Docker(ctx).SwarmInit(ctx, ir); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("%s: swarm initialize fail: %s", h.Id(), err.Error()))
			continue
		}

		// swarm init suceeded on at least one manager, so we have a leader
		slog.DebugContext(ctx, fmt.Sprintf("%s: swarm initialized. This is the new swarm leader", h.Id()), slog.Any("host", h))

		return h, nil
	}

	return nil, fmt.Errorf("could not initialize swarm on any manager host")
}

// join a host to the swarm
//
// @NOTE We have to check if the host is already in the swarm
//
//	and we may need to leave a/the swarm first
func joinSwarm(ctx context.Context, h *dockerhost.Host, sw dockertypesswarm.Swarm, swaas []string, role string) error {
	if si, err := h.Docker(ctx).SwarmInspect(ctx); err == nil {
		if si.ID != sw.ID {
			// host is already in a different swarm
			if lerr := h.Docker(ctx).SwarmLeave(ctx, true); lerr != nil {
				return fmt.Errorf("%s: was in another swarm '%s', and failed to leave it", h.Id(), si.ID)
			}

			// @TODO check it is not in the wrong role

		} else {
			slog.InfoContext(ctx, fmt.Sprintf("%s: host already in swarm", h.Id()), slog.Any("host", h))
			return nil // no action needed
		}
	}

	r := dockertypesswarm.JoinRequest{
		RemoteAddrs: swaas,
	}

	switch role {
	case "manager":
		r.JoinToken = sw.JoinTokens.Manager
	case "worker":
		r.JoinToken = sw.JoinTokens.Worker
	}

	if err := h.Docker(ctx).SwarmJoin(ctx, r); err != nil {
		return fmt.Errorf("%s: failed to join manager to swarm: %s", h.Id(), err.Error())
	}

	return nil
}
