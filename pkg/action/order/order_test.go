package order

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validateOrder(t *testing.T, original, sorted Orderables) {
	t.Helper()
	require.Equal(t, len(original), len(sorted), "Sorted size mismatch")

	// Map labels to their position in the sorted list
	deliveryMap := make(map[string][]int)
	for i, o := range sorted {
		for _, d := range o.Delivers {
			deliveryMap[d] = append(deliveryMap[d], i)
		}
	}

	for i, o := range sorted {
		// Must come AFTER all providers of labels in Before
		for _, b := range o.Before {
			providers, ok := deliveryMap[b]
			if !ok {
				t.Errorf("Item %s requires label %s which is not delivered", o.Key, b)
				continue
			}
			for _, providerIdx := range providers {
				if providerIdx > i {
					t.Errorf("Constraint violation: %s must be after delivery of %s (provided by item at index %d, currently at index %d)", o.Key, b, providerIdx, i)
				}
			}
		}

		// Must come BEFORE all providers of labels in After
		for _, a := range o.After {
			providers, ok := deliveryMap[a]
			if !ok {
				continue // After labels are optional
			}
			for _, providerIdx := range providers {
				if providerIdx < i {
					t.Errorf("Constraint violation: %s must be before delivery of %s (provided by item at index %d, currently at index %d)", o.Key, a, providerIdx, i)
				}
			}
		}
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name        string
		input       Orderables
		expectError error
	}{
		{
			name: "Linear dependency",
			input: Orderables{
				{Key: "step2", Before: []string{"a"}},
				{Key: "step1", Delivers: []string{"a"}},
			},
		},
		{
			name: "Complex graph",
			input: Orderables{
				{Key: "infrastructure", Delivers: []string{"cloud"}},
				{Key: "network", Delivers: []string{"vpc"}, Before: []string{"cloud"}},
				{Key: "compute", Delivers: []string{"nodes"}, Before: []string{"vpc"}},
				{Key: "storage", Delivers: []string{"ebs"}, Before: []string{"cloud"}},
				{Key: "database", Delivers: []string{"rds"}, Before: []string{"vpc", "ebs"}},
				{Key: "application", Delivers: []string{"web"}, Before: []string{"nodes", "rds"}},
				{Key: "monitoring", Before: []string{"web"}},
				{Key: "cleanup", After: []string{"web"}}, // cleanup must be before web
			},
		},
		{
			name: "Circular dependency",
			input: Orderables{
				{Key: "a", Delivers: []string{"label-a"}, Before: []string{"label-b"}},
				{Key: "b", Delivers: []string{"label-b"}, Before: []string{"label-a"}},
			},
			expectError: ErrCouldNotSort,
		},
		{
			name: "Missing dependency",
			input: Orderables{
				{Key: "a", Before: []string{"missing"}},
			},
			expectError: ErrSortDependencyNotDelivered,
		},
		{
			name: "Multiple delivers",
			input: Orderables{
				{Key: "p1", Delivers: []string{"a", "b"}},
				{Key: "c1", Before: []string{"a"}},
				{Key: "c2", Before: []string{"b"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted, err := Sort(tt.input)
			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
				return
			}

			require.NoError(t, err)
			validateOrder(t, tt.input, sorted)

			// Simple visual check for complex graph
			if tt.name == "Complex graph" {
				fmt.Printf("\nSorted order for '%s':\n", tt.name)
				for i, o := range sorted {
					fmt.Printf("%d: %s\n", i, o.Key)
				}
			}
		})
	}
}
