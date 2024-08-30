package dependency_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_NewDependencies(t *testing.T) {
	d0 := mock.StaticDependency("zero", "zero", nil, nil)
	d1 := mock.StaticDependency("one", "one", nil, nil)
	d2 := mock.StaticDependency("two", "two", nil, nil)
	d3 := mock.StaticDependency("three", "three", nil, nil)
	d4 := mock.StaticDependency("four", "four", nil, nil)

	ds := dependency.NewDependencies(d0, d1, d2, d3, d4)

	assert.Len(t, ds, 5, "New dependencies had wrong length")

	assert.Same(t, d0, ds[0], "wrong dependency in position")
	assert.Same(t, d1, ds[1], "wrong dependency in position")
	assert.Same(t, d2, ds[2], "wrong dependency in position")
	assert.Same(t, d3, ds[3], "wrong dependency in position")
	assert.Same(t, d4, ds[4], "wrong dependency in position")
}
