package action

import (
	"github.com/Mirantis/launchpad/pkg/dependency"
)

type Steps []Step

func (ss Steps) Ids() []string {
	ids := []string{}
	for _, id := range ids {
		ids = append(ids, id)
	}
	return ids
}

func OrderSteps(sps *Steps, reverse bool) error {
	os := orderables{}

	befores := map[string][]string{}

	for _, p := range *sps {
		pr := []string{}
		pp := []string{}

		pr = append(pr, p.Id) // provide the step

		pr = append(pr, p.Requires.Ids()...)
		pp = append(pp, p.Provides.Ids()...)

		pr = append(pr, p.After.Ids()...)

		for _, a := range p.Before.Ids() {
			if _, ok := befores[a]; !ok {
				befores[a] = []string{}
			}
			befores[a] = append(befores[a], p.Id)
		}

		os = append(os, orderable{
			id:       p.Id,
			requires: pr,
			provides: pp,
		})
	}

	for i, o := range os {
		b, ok := befores[o.id]
		if !ok {
			continue
		}

		o.requires = append(o.requires, b...)
		os[i] = o
	}

	oso, err := os.order(reverse)
	if err != nil {
		return err
	}

	ops := Steps{}
	for _, oi := range oso {
		p := (*sps)[oi]
		ops = append(ops, p)
	}
	*sps = ops

	return nil
}

type Step struct {
	Id string

	// Provides indiate that this action MAY provide a dependency
	Provides dependency.Dependencies
	Requires dependency.Dependencies

	After  Steps
	Before Steps
}
