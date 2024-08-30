package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComponentSanityPass(t *testing.T) {
	ctx := context.Background()
	cn := "compname"

	c := mock.NewSimpleComponent(cn, nil, nil)

	require.Equal(t, cn, c.Name(), "mock component had incorrect name: %s != %s", cn, c.Name())

	c_d := c.Debug()
	require.NotNil(t, c_d, "mock component gave non nil debug error: %s", c_d)
	c_v := c.Validate(ctx)
	require.Nil(t, c_v, "mock component gave non nil validate error: %+v", c_v)
}

func Test_ComponentSanityFail(t *testing.T) {
	ctx := context.Background()
	cv := errors.New("expected validate failure")

	c := mock.NewSimpleComponent("name", nil, cv)

	c_v := c.Validate(ctx)
	require.Error(t, c_v, "mock component did not give expected validate error: %+v", cv, c_v)
}

func Test_ComponentSanityDebug(t *testing.T) {
	cn := "compname"
	cd := map[string]string{"test": "test"}

	c := mock.NewSimpleComponent(cn, cd, nil)

	c_d := c.Debug()

	c_dd, ok := c_d.(mock.ComponentDebug)
	assert.True(t, ok, "debug returned unexpected type: %+v", c_dd)
	assert.Equal(t, c_dd.Name, cn, "incorrect debug name returned")
	assert.Equal(t, c_dd.Debug, cd, "incorrect debug returned")
}

func Test_ComponentNilFactories(t *testing.T) {
	ctx := context.Background()

	c := mock.NewComponentWDependencies("my-component", nil, nil, nil, nil)

	assert.Nil(t, c.RequiresDependencies(ctx), "wrong requirements returned")

	cd, cderr := c.ProvidesDependencies(ctx, nil)

	assert.Nil(t, cderr, "")
	assert.Nil(t, cd, "Non nil error returned")
}

func Test_ComponentNilRequiresDependencies(t *testing.T) {
	ctx := context.Background()

	fac := func(_ context.Context) dependency.Requirements {
		return nil
	}

	c := mock.NewComponentWDependencies("my-component", nil, nil, fac, nil)

	assert.Empty(t, c.RequiresDependencies(ctx), "wrong requirements returned")
}

func Test_ComponentRequiresDependencies(t *testing.T) {
	ctx := context.Background()

	rs := dependency.NewRequirements()
	fac := func(_ context.Context) dependency.Requirements {
		return rs
	}

	c := mock.NewComponentWDependencies("my-component", nil, nil, fac, nil)

	assert.Equal(t, rs, c.RequiresDependencies(ctx), "wrong requirements returned")
}

func Test_ComponentNilProvidesDependencies(t *testing.T) {
	ctx := context.Background()

	err := errors.New("expected factory error")
	fac := func(_ context.Context, _ dependency.Requirement) (dependency.Dependency, error) {
		return nil, err
	}

	c := mock.NewComponentWDependencies("my-component", nil, nil, nil, fac)

	cd, cderr := c.ProvidesDependencies(ctx, nil)

	assert.ErrorIs(t, err, cderr, "")
	assert.Nil(t, cd, "Non nil error returned")
}

func Test_ComponentProvidesDependencies(t *testing.T) {
	ctx := context.Background()

	d := mock.SimpleDependency("my-dep", "my test dependency", nil, nil)
	fac := func(_ context.Context, _ dependency.Requirement) (dependency.Dependency, error) {
		return d, nil
	}

	c := mock.NewComponentWDependencies("my-component", nil, nil, nil, fac)

	cd, cderr := c.ProvidesDependencies(ctx, nil)

	assert.Nil(t, cderr, "Non nil error returned")
	assert.Same(t, d, cd, "wrong requirements returned")
}
