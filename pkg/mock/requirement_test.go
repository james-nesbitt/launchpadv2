package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_RequirementSanity(t *testing.T) {
	d := mock.StaticDependency("test-dep", "my simple dep", nil, nil)

	assert.Implements(t, (*dependency.Requirement)(nil), mock.StaticRequirement("test-req", "my simple req", d), "Mock req does not implement expected dependency.Requirement interface")
}

func Test_RequirementStatic(t *testing.T) {
	ctx := context.Background()

	rid := "test-req"
	rdesc := "my test req"
	d := mock.StaticDependency("test-dep", "my simple dep", nil, nil)

	r := mock.StaticRequirement(rid, rdesc, d)

	assert.Equal(t, rid, r.Id(), "wrong id returned")
	assert.Equal(t, rdesc, r.Describe(), "wrong description returned")

	rd := r.Matched(ctx)
	assert.Same(t, d, rd, "wrong dependency returned")

	// this req type is supposed to not allow new dependencies to be assigned.
	assert.Error(t, r.Match(nil), "static requirement did not return expected error when given a new dependency")
}

func Test_RequirementSimple(t *testing.T) {
	ctx := context.Background()

	rid := "simple-req"
	rdesc := "simple requirement"

	r := mock.SimpleRequirement(rid, rdesc)

	assert.Nil(t, r.Matched(ctx), "new simple requirement gave non-nil dependency, before being given one")

	d1 := mock.StaticDependency("test-dep-1", "my simple dep", nil, nil)
	d2 := mock.StaticDependency("test-dep-2", "my simple dep", nil, nil)

	assert.Nil(t, r.Match(d1), "simple requirement gave non-nil error when given first dependency")
	assert.Equal(t, d1, r.Matched(ctx), "simple requirement did not return expected first dependency after matching")

	assert.Nil(t, r.Match(d2), "simple requirement gave non-nil error when given second dependency")
	assert.Equal(t, d2, r.Matched(ctx), "simple requirement did not return expected second dependency after matching")
}

func Test_RequirementMatch(t *testing.T) {
	ctx := context.Background()
	d := mock.StaticDependency("test-dep", "my simple dep", nil, nil)

	rpos := mock.MatchingRequirement("matcher", "matching-requirement", func(_ dependency.Dependency) error { return nil })

	assert.Nil(t, rpos.Match(d), "matching requirement returned unexpected error on match")
	assert.Equal(t, d, rpos.Matched(ctx), "matching requirement did not return expected dependency after matching")

	rneg := mock.MatchingRequirement("matcher", "non-matching requirement", func(_ dependency.Dependency) error { return errors.New("expected matching error") })

	assert.Error(t, rneg.Match(d), "non matching requirement did not give expected error when matching")
	assert.Nil(t, rneg.Matched(ctx), "non matching requirement returned a matched after expected matching failure")
}
