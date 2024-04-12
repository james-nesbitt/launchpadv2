package action // testing internal functionality

import "testing"

func Test_OrderFails(t *testing.T) {
	os := orderables{
		orderable{
			requires: []string{"A"},
			provides: []string{"A"},
		},
	}

	_, err := os.order(false)
	if err == nil {
		t.Error("expected error with undeliveable order")
	}
}

func Test_Order(t *testing.T) {
	os := orderables{
		orderable{
			requires: []string{"A"},
			provides: []string{"E"},
		},
		orderable{
			requires: []string{"E"},
			provides: []string{"D"},
		},
		orderable{
			requires: []string{},
			provides: []string{"A"},
		},
		orderable{
			requires: []string{"B"},
			provides: []string{"C"},
		},
		orderable{
			requires: []string{},
			provides: []string{"B"},
		},
		orderable{
			requires: []string{"B"},
			provides: []string{},
		},
	}

	oso, err := os.order(false)
	if err != nil {
		t.Fatalf("order error unexpected: %s", err.Error())
	}

	if len(oso) != len(os) {
		t.Error("order gave incorrect number of entries")
	}

	if oso[0] != 2 {
		t.Errorf("Bad order 0")
	}
	if oso[5] != 1 {
		t.Errorf("Bad order 5")
	}
}
