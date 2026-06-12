package order_test

import (
	"errors"
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

func TestSort_ComplexCases(t *testing.T) {
	tests := []struct {
		name    string
		os      order.Orderables
		want    []string
		wantErr error
	}{
		{
			name: "Linear Chain",
			os: order.Orderables{
				{Key: "C", Delivers: []string{"C"}, Before: []string{"B"}},
				{Key: "A", Delivers: []string{"A"}},
				{Key: "B", Delivers: []string{"B"}, Before: []string{"A"}},
			},
			want: []string{"A", "B", "C"},
		},
		{
			name: "Disconnected Components",
			os: order.Orderables{
				{Key: "A", Delivers: []string{"A"}},
				{Key: "B", Delivers: []string{"B"}},
			},
			want: []string{"A", "B"},
		},
		{
			name: "Cycle Detection",
			os: order.Orderables{
				{Key: "A", Delivers: []string{"A"}, Before: []string{"B"}},
				{Key: "B", Delivers: []string{"B"}, Before: []string{"A"}},
			},
			wantErr: order.ErrCouldNotSort,
		},
		{
			name: "Missing Dependency",
			os: order.Orderables{
				{Key: "A", Before: []string{"Missing"}},
			},
			wantErr: order.ErrSortDependencyNotDelivered,
		},
		{
			name: "Empty List",
			os:   order.Orderables{},
			want: []string{},
		},
		{
			name: "Single Element",
			os: order.Orderables{
				{Key: "A", Delivers: []string{"A"}},
			},
			want: []string{"A"},
		},
		{
			name: "Complex DAG",
			os: order.Orderables{
				{Key: "1", Delivers: []string{"D1"}},
				{Key: "2", Delivers: []string{"D2"}},
				{Key: "3", Before: []string{"D1", "D2"}},
				{Key: "4", Before: []string{"D3"}},
				{Key: "5", Delivers: []string{"D3"}, Before: []string{"D1"}},
				{Key: "6", Before: []string{"D2", "D3"}},
			},
			want: []string{"1", "2", "5", "3", "4", "6"}, 
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.Sort(tt.os)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Sort() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Sort() unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("Sort() length = %v, want %v", len(got), len(tt.want))
				return
			}
			
			if err := verifyOrder(tt.os, got); err != nil {
				t.Errorf("Sort() produced invalid order: %v", err)
			}
		})
	}
}

func verifyOrder(os order.Orderables, sorted order.Orderables) error {
	if len(os) != len(sorted) {
		return fmt.Errorf("length mismatch: os=%d, sorted=%d", len(os), len(sorted))
	}

	pos := make(map[string]int)
	for i, o := range sorted {
		pos[o.Key] = i
	}

	for _, o := range os {
		myPos := pos[o.Key]

		for _, label := range o.Before {
			foundProvider := false
			for _, provider := range os {
				for _, dl := range provider.Delivers {
					if dl == label {
						foundProvider = true
						if pos[provider.Key] >= myPos {
							return fmt.Errorf("constraint violation: provider %s of %s must come before %s", provider.Key, label, o.Key)
						}
					}
				}
			}
			if !foundProvider {
				return fmt.Errorf("missing provider for label %s required by %s", label, o.Key)
			}
		}

		for _, label := range o.After {
			foundProvider := false
			for _, provider := range os {
				for _, dl := range provider.Delivers {
					if dl == label {
						foundProvider = true
						if pos[provider.Key] <= myPos {
							return fmt.Errorf("constraint violation: provider %s of %s must come after %s", provider.Key, label, o.Key)
						}
					}
				}
			}
			if !foundProvider {
				return fmt.Errorf("missing provider for label %s required by %s", label, o.Key)
			}
		}
	}
	return nil
}
func TestSort_StressTest(t *testing.T) {
	const count = 1000
	os := order.Orderables{}

	for i := 0; i < count; i++ {
		o := order.Orderable{
			Key:      fmt.Sprintf("%d", i),
			Delivers: []string{fmt.Sprintf("D%d", i)},
		}
		// Each element (except the first) depends on the one before it
		if i > 0 {
			o.Before = append(o.Before, fmt.Sprintf("D%d", i-1))
		}
		// Add some random breadth constraints to increase complexity
		if i > 10 && i%10 == 0 {
			o.Before = append(o.Before, fmt.Sprintf("D%d", i-10))
		}
		os = append(os, o)
	}

	sos, err := order.Sort(os)
	if err != nil {
		t.Fatalf("Stress test sort unexpected error: %v", err)
	}

	if len(sos) != count {
		t.Errorf("Stress test length mismatch: got %d, want %d", len(sos), count)
	}

	if err := verifyOrder(os, sos); err != nil {
		t.Errorf("Stress test produced invalid order: %v", err)
	}
}
