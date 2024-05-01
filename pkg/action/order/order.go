package order

import (
	"errors"
	"fmt"

	"github.com/yourbasic/graph"
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

	g := graph.New(len(os))

	// for any label, add an aedge between the delivers and the after/before elements

	rerrs := []error{}
	for k, l := range ls {
		db := ls[k].deliveredBy

		if len(l.before) > 0 && len(db) == 0 {
			rerrs = append(rerrs, fmt.Errorf("%s is not delivered, but is required", k))
		}

		// things that come "before" a "delivers" get a pos *s cost
		for _, b := range l.before {
			for i, d := range db {
				g.AddCost(d, b, int64(i*2))
			}
		}
		// things that come "after" a "delivers" get a pos cost
		for _, a := range l.after {
			for i, d := range db {
				g.AddCost(a, d, int64(i))
			}
		}
	}

	if len(rerrs) > 0 {
		return Orderables{}, fmt.Errorf("%w; %s", ErrSortDependencyNotDelivered, errors.Join(rerrs...).Error())
	}

	soi, ok := graph.TopSort(g)
	if !ok {
		return Orderables{}, ErrCouldNotSort
	}

	sos := Orderables{}
	for _, oi := range soi {
		sos = append(sos, os[oi])
	}

	return sos, nil
}
