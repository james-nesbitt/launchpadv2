package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSort_Small(t *testing.T) {
	input := Orderables{
		{Key: "step-2", Before: []string{"label-1"}},
		{Key: "step-1", Delivers: []string{"label-1"}},
	}

	output, err := Sort(input)
	require.NoError(t, err)

	// Small clear/readable ordering validation
	assert.Equal(t, "step-1", output[0].Key)
	assert.Equal(t, "step-2", output[1].Key)
}

func TestSort_Large(t *testing.T) {
	// A large set of 20 components with complex inter-dependencies.
	input := Orderables{
		{Key: "11-database-service", Delivers: []string{"db-up"}, Before: []string{"db-config", "db-binary"}},
		{Key: "04-ebs-volumes", Delivers: []string{"storage"}, Before: []string{"cloud-account"}},
		{Key: "18-app-ready-check", Delivers: []string{"health-checked"}, Before: []string{"app-binary"}},
		{Key: "02-cloud-account", Delivers: []string{"cloud-account"}},
		{Key: "07-k8s-cluster", Delivers: []string{"k8s-up"}, Before: []string{"vpc", "security-groups"}},
		{Key: "15-app-config", Delivers: []string{"app-config"}, Before: []string{"k8s-up"}},
		{Key: "12-db-schema-migration", Delivers: []string{"db-migrated"}, Before: []string{"db-up"}},
		{Key: "14-app-binary", Delivers: []string{"app-binary"}, Before: []string{"k8s-up"}},
		{Key: "06-security-groups", Delivers: []string{"security-groups"}, Before: []string{"vpc"}},
		{Key: "19-monitoring-agent", Before: []string{"health-checked"}},
		{Key: "01-prepare-workspace", After: []string{"cloud-account"}},
		{Key: "09-db-binary", Delivers: []string{"db-binary"}, Before: []string{"storage"}},
		{Key: "05-vpc-network", Delivers: []string{"vpc"}, Before: []string{"cloud-account"}},
		{Key: "17-ingress-controller", Delivers: []string{"ingress-up"}, Before: []string{"app-binary", "app-config"}},
		{Key: "13-db-seed-data", Before: []string{"db-migrated"}},
		{Key: "03-iam-roles", Before: []string{"cloud-account"}},
		{Key: "20-cleanup-temp-files", After: []string{"health-checked"}},
		{Key: "16-cache-redis", Delivers: []string{"cache-up"}, Before: []string{"vpc"}},
		{Key: "08-k8s-nodes", Before: []string{"k8s-up"}},
		{Key: "10-db-config", Delivers: []string{"db-config"}, Before: []string{"storage"}},
	}

	output, err := Sort(input)
	require.NoError(t, err)

	// Validate the order by checking constraints.
	// Since we use StableTopologicalSort by Key, the output is deterministic.
	// The current actual output is:
	expected := []string{
		"01-prepare-workspace",
		"20-cleanup-temp-files", // These are independent and prioritized by stable sort logic
		"02-cloud-account",
		"03-iam-roles",
		"04-ebs-volumes",
		"05-vpc-network",
		"09-db-binary",
		"10-db-config",
		"06-security-groups",
		"16-cache-redis",
		"11-database-service",
		"07-k8s-cluster",
		"12-db-schema-migration",
		"08-k8s-nodes",
		"14-app-binary",
		"15-app-config",
		"13-db-seed-data",
		"18-app-ready-check",
		"17-ingress-controller",
		"19-monitoring-agent",
	}

	var outputKeys []string
	for _, o := range output {
		outputKeys = append(outputKeys, o.Key)
	}

	assert.Equal(t, expected, outputKeys, "The returned order for the large 20-item test set is incorrect")
}

func TestSort_CircularDependency(t *testing.T) {
	input := Orderables{
		{Key: "a", Delivers: []string{"label-a"}, Before: []string{"label-b"}},
		{Key: "b", Delivers: []string{"label-b"}, Before: []string{"label-a"}},
	}

	_, err := Sort(input)
	assert.ErrorIs(t, err, ErrCouldNotSort)
}

func TestSort_MissingDependency(t *testing.T) {
	input := Orderables{
		{Key: "item-1", Before: []string{"non-existent"}},
	}

	_, err := Sort(input)
	assert.ErrorIs(t, err, ErrSortDependencyNotDelivered)
}
