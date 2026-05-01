/*
Package order ordering functionality

Primarily used to order dependencies and steps.
*/
package order

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"
)

var (
	ErrSortDependencyNotDelivered = errors.New("sorting dependency not delivered")
	ErrCouldNotSort               = errors.New("could not sort")
)

type Orderables []Orderable

type Orderable struct {
	Key      string   // element identifier
	Delivers []string // all the labels which this element delivers
	Before   []string // all the labels which must go before this element
	After    []string // all the labels which must go after this element (IF THEY ARE FOUND)
}

type labels map[string]*label

func (ls labels) ensure(k string) {
	if _, ok := ls[k]; ok {
		return
	}

	ls[k] = &label{}
}

type label struct {
	deliveredBy []int
	before      []int
	after       []int
}

// Sort orderables by returning the index reorder.
func Sort(os Orderables) (Orderables, error) {
	ls := labels{}

	// collect all of the relational data from the orderable labels to the orderables
	for i, o := range os {
		for _, d := range o.Delivers {
			ls.ensure(d)
			ls[d].deliveredBy = append(ls[d].deliveredBy, i)
		}
		for _, d := range o.Before {
			ls.ensure(d)
			ls[d].before = append(ls[d].before, i)
		}
		for _, d := range o.After {
			ls.ensure(d)
			ls[d].after = append(ls[d].after, i)
		}
	}

	g := graph.New(graph.IntHash, graph.Directed())

	for i := range os {
		if err := g.AddVertex(i); err != nil {
			return nil, fmt.Errorf("failed to add vertex: %w", err)
		}
	}

	rerrs := []error{}
	for k, l := range ls {
		db := l.deliveredBy

		if len(l.before) > 0 && len(db) == 0 {
			rerrs = append(rerrs, fmt.Errorf("%s is not delivered, but is required", k))
		}

		// things that come "before" this element (labels)
		// so: deliverer -> element
		for _, b := range l.before {
			for _, d := range db {
				if err := g.AddEdge(d, b); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, fmt.Errorf("failed to add edge (before): %w", err)
				}
			}
		}
		// things that come "after" this element (labels)
		// so: element -> deliverer
		for _, a := range l.after {
			for _, d := range db {
				if err := g.AddEdge(a, d); err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
					return nil, fmt.Errorf("failed to add edge (after): %w", err)
				}
			}
		}
	}

	if len(rerrs) > 0 {
		return Orderables{}, fmt.Errorf("%w; %s", ErrSortDependencyNotDelivered, errors.Join(rerrs...).Error())
	}

	soi, err := graph.StableTopologicalSort(g, func(a, b int) bool {
		return os[a].Key < os[b].Key
	})
	if err != nil {
		return Orderables{}, fmt.Errorf("%w: %v", ErrCouldNotSort, err)
	}

	sos := Orderables{}
	for _, oi := range soi {
		sos = append(sos, os[oi])
	}

	return sos, nil
}
