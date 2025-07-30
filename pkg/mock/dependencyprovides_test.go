package mock_test

import (
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_ProvidesNilDependencies(t *testing.T) {
	ctx := t.Context()

	err := errors.New("expected error")
	pd := mock.ProvidesDependencies(nil, err)

	assert.Implements(t, (*dependency.ProvidesDependencies)(nil), pd, "mock provides-dependencies does not implements the dependency.ProvidesRequirements interface")

	r := mock.SimpleRequirement("test-req", "test requirement")
	pdd, derr := pd.ProvidesDependencies(ctx, r)

	assert.ErrorIs(t, derr, err, "did not receive expected error")
	assert.Nil(t, pdd, "mock requires-dependencies did not return expected value")
}

func Test_ProvidesDependencies(t *testing.T) {
	ctx := t.Context()

	d := mock.SimpleDependency("test-dep", "test dependency", nil, nil)
	pd := mock.ProvidesDependencies(d, nil)

	assert.Implements(t, (*dependency.ProvidesDependencies)(nil), pd, "mock provides-dependencies does not implements the dependency.RequiresDependencies interface")

	r := mock.SimpleRequirement("test-req", "test requirement")
	pdd, derr := pd.ProvidesDependencies(ctx, r)

	assert.Nil(t, derr, "mock provides-dependencies gave unexpected error")
	assert.Same(t, d, pdd, "mock provides-dependencies did not give the expected dependency")
}
