package mcr

import (
	"context"
	"fmt"
	"log/slog"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

type swarmActivateStep struct {
	baseStep

	id string

	managerJoinToken string
	workerJoinToken  string
}

func (s swarmActivateStep) Id() string {
	return fmt.Sprintf("%s:mcr-swarm-activate", s.id)
}

func (s swarmActivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR swarm-activate step", slog.String("ID", s.Id()))

	mhs, mhsgerr := s.c.GetManagerHosts(ctx)
	if mhsgerr != nil {
		return mhsgerr
	}

	if s.c.state.leader == nil {
		for _, mh := range mhs {
			if sierr := s.initSwarm(ctx, mh); sierr != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("%s: swarm initialize fail: %s", mh.Id(), sierr.Error()))
				continue
			}

			// swarm init suceeded on at least one machine
			s.c.state.leader = mh
			break
		}

		if s.c.state.leader == nil {
			return fmt.Errorf("Swarm init failed on all Managers")
		}
	}

	if err := mhs.Each(ctx, s.joinSwarmManager); err != nil {
		return err
	}

	whs, whsgerr := s.c.GetWorkerHosts(ctx)
	if whsgerr != nil {
		return whsgerr
	}
	if err := whs.Each(ctx, s.joinSwarmManager); err != nil {
		return err
	}

	return nil
}

// initSwarm Initialize swarm on the host
//
// @NOTE Run this only on one host
func (s swarmActivateStep) initSwarm(ctx context.Context, h *dockerhost.Host) error {
	n, nerr := h.Network(ctx)
	if nerr != nil {
		return nerr
	}

	sf := []string{}

	sf = append(sf, fmt.Sprintf("--advertise-addr=%s", n.PrivateAddress))
	sf = append(sf, s.c.config.SwarmInstallFlags...)

	h.Docker(ctx).SwarmInit(ctx, sf)

	s.c.state.leader = h
	return nil
}

func (s swarmActivateStep) joinSwarmManager(ctx context.Context, h *dockerhost.Host) error {
	if s.c.state.leader.Id() == h.Id() { // the leader is already in the swarm
		return nil
	}

	l := s.c.state.leader
	ln, lnerr := l.Network(ctx)
	if lnerr != nil {
		return fmt.Errorf("no leader network, can't join its swarm: %s", lnerr)
	}

	// @TODO This should be just collected once and put into state
	t, terr := l.Docker(ctx).SwarmJoinToken(ctx, "manager", false)
	if terr != nil {
		return fmt.Errorf("no leader token, can't join its swarm: %s", terr)
	}

	if err := h.Docker(ctx).SwarmJoin(ctx, ln.PrivateAddress, t, []string{}); err != nil {
		return fmt.Errorf("failed to join swarm: %s", err.Error())
	}

	return nil
}

func (s swarmActivateStep) joinSwarmWorker(ctx context.Context, h *dockerhost.Host) error {
	l := s.c.state.leader

	ln, lnerr := l.Network(ctx)
	if lnerr != nil {
		return fmt.Errorf("no leader network, can't join its swarm: %s", lnerr)
	}

	// @TODO This should be just collected once and put into state
	t, terr := l.Docker(ctx).SwarmJoinToken(ctx, "worker", false)
	if terr != nil {
		return fmt.Errorf("no leader token, can't join its swarm: %s", terr)
	}

	if err := h.Docker(ctx).SwarmJoin(ctx, ln.PrivateAddress, t, []string{}); err != nil {
		return fmt.Errorf("failed to join swarm: %s", err.Error())
	}

	return nil
}
