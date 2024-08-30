package mock_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_RequiresNilDependencies(t *testing.T) {
	ctx := context.Background()

	rd := mock.RequiresDependencies(nil)

	assert.Implements(t, (*dependency.RequiresDependencies)(nil), rd, "mock requires-dependencies does not implements the dependency.RequiresDependencies interface")
	assert.Nil(t, rd.RequiresDependencies(ctx), "mock requires-dependencies did not return expected value")
}

func Test_RequiresDependencies(t *testing.T) {
	ctx := context.Background()

	rs := dependency.NewRequirements()
	rd := mock.RequiresDependencies(rs)

	assert.Implements(t, (*dependency.RequiresDependencies)(nil), rd, "mock requires-dependencies does not implements the dependency.RequiresDependencies interface")
	assert.Equal(t, rs, rd.RequiresDependencies(ctx), "mock requires-dependencies did not return expected value")
}
