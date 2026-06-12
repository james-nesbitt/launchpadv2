/*
Package order ordering functionality

Primarily used to order dependencies and steps.
*/
package order

import (
	"errors"
	"fmt"

	// Removed yourbasic/graph
)

var (
	ErrSortDependencyNotDelivered = errors.New("sorting dependecy not delivered")
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

	edges := [][2]int{}
	rerrs := []error{}
	for k, l := range ls {
		db := l.deliveredBy

		if len(l.before) > 0 && len(db) == 0 {
			rerrs = append(rerrs, fmt.Errorf("%s is not delivered, but is required", k))
		}

		// things that come "before" a "delivers" get an edge
		for _, b := range l.before {
			for _, d := range db {
				edges = append(edges, [2]int{d, b})
			}
		}
		// things that come "after" a "delivers" get an edge
		for _, a := range l.after {
			for _, d := range db {
				edges = append(edges, [2]int{a, d})
			}
		}
	}

	if len(rerrs) > 0 {
		return Orderables{}, fmt.Errorf("%w; %s", ErrSortDependencyNotDelivered, errors.Join(rerrs...).Error())
	}

	soi, ok := topoSort(len(os), edges)
	if !ok {
		return Orderables{}, ErrCouldNotSort
	}

	sos := Orderables{}
	for _, oi := range soi {
		sos = append(sos, os[oi])
	}

	return sos, nil
}

func topoSort(nodes int, edges [][2]int) ([]int, bool) {
	adj := make([][]int, nodes)
	inDegree := make([]int, nodes)
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		adj[u] = append(adj[u], v)
		inDegree[v]++
	}

	var queue []int
	for i := 0; i < nodes; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	var result []int
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)

		for _, v := range adj[u] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(result) != nodes {
		return nil, false
	}
	return result, true
}
