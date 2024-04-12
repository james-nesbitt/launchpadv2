package action

import (
	"fmt"

	"github.com/yourbasic/graph"
)

type orderable struct {
	id       string
	requires []string
	provides []string
}

type orderables []orderable

func (os orderables) order(reverse bool) ([]int, error) {
	m := map[string]struct {
		r []int // required by
		p []int // provided by
	}{}
	ensureInMap := func(id string) {
		if _, ok := m[id]; !ok {
			m[id] = struct {
				r []int // required by
				p []int // provided by
			}{
				r: []int{},
				p: []int{},
			}
		}
	}

	for i, o := range os {
		for _, id := range o.provides {
			ensureInMap(id)

			mid := m[id]
			mid.p = append(mid.p, i)
			m[id] = mid
		}
		for _, id := range o.requires {
			ensureInMap(id)

			mid := m[id]
			mid.r = append(mid.r, i)
			m[id] = mid
		}
	}

	pg := graph.New(len(os))

	for _, d := range m {
		for _, p := range d.p {
			for _, r := range d.r {
				pg.Add(p, r)
			}
		}
	}

	var o []int
	var ok bool

	if reverse {
		rpg := graph.Transpose(pg)
		o, ok = graph.TopSort(rpg)
	} else {
		o, ok = graph.TopSort(pg)
	}

	if !ok {
		return o, fmt.Errorf("could not re-order: %+v", pg)
	}

	return o, nil
}
