package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_DependencySanity(t *testing.T) {
	assert.Implements(t, (*dependency.Dependency)(nil), mock.StaticDependency("test-dep", "my simple dep", nil, nil), "Mock dep does not implement expected dependency.Dependency interface")
}

func Test_DependencyStaticNoErrors(t *testing.T) {
	ctx := t.Context()

	id := "my-dep"
	desc := "my dependency"

	d := mock.StaticDependency(id, desc, nil, nil)

	assert.Equal(t, id, d.Id(), "mock dep returned wrong ID")
	assert.Equal(t, desc, d.Describe(), "mock dep returned wrong description")

	assert.Nil(t, d.Validate(ctx), "mock dependency did not return expected nil validation")
	assert.Nil(t, d.Met(ctx), "mock dependency did not return expected nil Met")
}

func Test_DependencyStaticWErrors(t *testing.T) {
	ctx := t.Context()

	id := "my-dep"
	desc := "my dependency"
	verr := errors.New("expected validation error")
	merr := errors.New("expected met error")

	d := mock.StaticDependency(id, desc, verr, merr)

	assert.Equal(t, id, d.Id(), "mock dep returned wrong ID")
	assert.Equal(t, desc, d.Describe(), "mock dep returned wrong description")

	assert.ErrorIs(t, d.Validate(ctx), verr, "mock dependency did not return expected error")
	assert.ErrorIs(t, d.Met(ctx), merr, "mock dependency did not return expected Met error")
}

func Test_DependencySimpleNil(t *testing.T) {
	ctx := t.Context()

	id := "my-dep"
	desc := "my dependency"

	d := mock.SimpleDependency(id, desc, func(_ context.Context) error { return nil }, func(_ context.Context) error { return nil })

	assert.Nil(t, d.Validate(ctx), "mock dependency did not return expected nil validation")
	assert.Nil(t, d.Met(ctx), "mock dependency did not return expected nil Met")
}

func Test_DependencySimple(t *testing.T) {
	ctx := t.Context()

	id := "my-dep"
	desc := "my dependency"
	verr := errors.New("expected validation error")
	merr := errors.New("expected met error")

	d := mock.SimpleDependency(id, desc, func(_ context.Context) error { return verr }, func(_ context.Context) error { return merr })

	assert.ErrorIs(t, d.Validate(ctx), verr, "mock dependency did not return expected error")
	assert.ErrorIs(t, d.Met(ctx), merr, "mock dependency did not return expected Met error")
}
