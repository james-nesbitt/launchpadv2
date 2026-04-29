package k0s

import "testing"

func Test_mergeUnique(t *testing.T) {
	tests := []struct {
		name      string
		base      []string
		additions []string
		want      []string
	}{
		{
			name:      "empty inputs",
			base:      nil,
			additions: nil,
			want:      []string{},
		},
		{
			name:      "additions only",
			base:      nil,
			additions: []string{"a", "b"},
			want:      []string{"a", "b"},
		},
		{
			name:      "base only",
			base:      []string{"a", "b"},
			additions: nil,
			want:      []string{"a", "b"},
		},
		{
			name:      "deduplication",
			base:      []string{"a", "b"},
			additions: []string{"b", "c"},
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "all duplicates",
			base:      []string{"a", "b"},
			additions: []string{"a", "b"},
			want:      []string{"a", "b"},
		},
		{
			name:      "order preserved (base first)",
			base:      []string{"z"},
			additions: []string{"a", "b"},
			want:      []string{"z", "a", "b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeUnique(tc.base, tc.additions)
			if len(got) != len(tc.want) {
				t.Fatalf("mergeUnique() len=%d, want %d: got=%v want=%v", len(got), len(tc.want), got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("mergeUnique()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func Test_digStringSlice(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]interface{}
		key  string
		want []string
	}{
		{
			name: "missing key",
			m:    map[string]interface{}{},
			key:  "sans",
			want: nil,
		},
		{
			name: "nil value",
			m:    map[string]interface{}{"sans": nil},
			key:  "sans",
			want: nil,
		},
		{
			name: "wrong type",
			m:    map[string]interface{}{"sans": "not a slice"},
			key:  "sans",
			want: nil,
		},
		{
			name: "string slice via interface",
			m:    map[string]interface{}{"sans": []interface{}{"a.example.com", "b.example.com"}},
			key:  "sans",
			want: []string{"a.example.com", "b.example.com"},
		},
		{
			name: "mixed types skips non-strings",
			m:    map[string]interface{}{"sans": []interface{}{"a", 42, "b"}},
			key:  "sans",
			want: []string{"a", "b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := digStringSlice(tc.m, tc.key)
			if len(got) != len(tc.want) {
				t.Fatalf("digStringSlice() = %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("digStringSlice()[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}
