package action

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/action/order"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	// Command phase which is used to initialize any components, usually requireing no
	//   interactions with a project
	CommandPhaseInitialize = "initialize"
	// Command phase which is used to discover state.
	CommandPhaseDiscover = "discover"
	// Command phase which is used to finalize and cleanup.
	CommandPhaseFinalize = "finalize"
)

type Phases map[string]Phase

func NewPhases(psa ...Phase) Phases {
	ps := Phases{}
	for _, p := range psa {
		ps.Add(p)
	}
	return ps
}

func (ps Phases) Contains(id string) bool {
	_, ok := ps[id]
	return ok
}

func (ps Phases) Add(nps ...Phase) {
	for _, np := range nps {
		// @TODO should we test for empty string?
		ps[np.Id()] = np
	}
}

func (ps Phases) Merge(nps Phases) {
	for _, np := range nps {
		// @TODO should we test for empty string?
		ps[np.Id()] = np
	}
}

func (ps Phases) Get(id string) (Phase, bool) {
	p, ok := ps[id]
	return p, ok
}

func (ps Phases) Order(ctx context.Context) ([]Phase, error) {
	os := order.Orderables{}

	for _, p := range ps {
		d := []string{}
		if pds, ok := p.(dependency.DeliversEvents); ok {
			for _, pd := range pds.DeliversEvents(ctx) {
				d = append(d, pd.Id)
			}
		}

		b := []string{}
		a := []string{}
		if res, ok := p.(dependency.RequiresEvents); ok {
			bes, aes := res.RequiresEvents(ctx)
			for _, pd := range bes {
				b = append(b, pd.Id)
			}

			for _, pd := range aes {
				a = append(a, pd.Id)
			}
		}

		os = append(os, order.Orderable{
			Key:      p.Id(),
			Delivers: d,
			Before:   b,
			After:    a,
		})
	}

	ops := []Phase{}
	pso, err := order.Sort(os)
	if err != nil {
		return ops, err
	}

	for _, o := range pso {
		ops = append(ops, ps[o.Key])
	}

	return ops, nil
}
