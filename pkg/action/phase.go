package action

import (
	"github.com/Mirantis/launchpad/pkg/dependency"
)

type Phases []Phase

func OrderPhases(psp *Phases, reverse bool) error {
	os := orderables{}
	for _, p := range *psp {
		os = append(os, orderable{
			requires: p.Requires.Ids(),
			provides: p.Provides.Ids(),
		})
	}

	oso, err := os.order(reverse)
	if err != nil {
		return err
	}

	ops := Phases{}
	for _, oi := range oso {
		p := (*psp)[oi]
		ops = append(ops, p)
	}
	*psp = ops

	return nil
}

type Phase struct {
	Id string

	Provides dependency.Dependencies
	Requires dependency.Dependencies
}

func PhaseSteps(p *Phase) (Steps, error) {
	ss := Steps{}

	if err := OrderSteps(&ss, false); err != nil {
		return ss, err
	}

	return ss, nil
}
