package order_test

import (
	"fmt"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action/order"
)

func Test_Ordering(t *testing.T) {
	os := order.Orderables{}

	os = append(os, order.Orderable{
		Key:      "2",
		Delivers: []string{},
		Before:   []string{"A"},
		After:    []string{},
	})
	os = append(os, order.Orderable{
		Key:      "3",
		Delivers: []string{"B"},
		Before:   []string{"A"},
		After:    []string{"D"},
	})
	os = append(os, order.Orderable{
		Key:      "0",
		Delivers: []string{"X"},
		Before:   []string{},
		After:    []string{"A"},
	})
	os = append(os, order.Orderable{
		Key:      "6",
		Delivers: []string{"E"},
		Before:   []string{"B", "C"},
		After:    []string{},
	})
	os = append(os, order.Orderable{
		Key:      "1",
		Delivers: []string{"A"},
		Before:   []string{},
		After:    []string{},
	})
	os = append(os, order.Orderable{
		Key:      "4",
		Delivers: []string{"C"},
		Before:   []string{"B"},
		After:    []string{"D"},
	})
	os = append(os, order.Orderable{
		Key:      "5",
		Delivers: []string{"D"},
		Before:   []string{},
		After:    []string{"E"},
	})

	if len(os) != 7 {
		t.Errorf("orderables length is wrong: %+v", os)
	}

	sos, err := order.Sort(os)
	if err != nil {
		t.Errorf("orderables sort unexpected error: %s", err.Error())
	}
	if len(sos) != len(os) {
		t.Errorf("not enough orderables returned in sort operation: %+v - %+v", os, sos)
	}

	for i, so := range sos {
		if fmt.Sprintf("%d", i) != so.Key {
			t.Errorf("orderable in wrong order [%d] %+v", i, so)
		} else {
			t.Logf("orderable in right order [%d] %+v", i, so)
		}
	}
}
