package mock_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_MockDependencyReqStatic(t *testing.T) {
	ctx := context.Background()

	rnfn := mock.ReqFactoryStatic(nil)
	rsn := rnfn(ctx)

	assert.Nil(t, rsn, "Nil Req Static function was supposed to return a nil")

	rss := dependency.Requirements{}
	rnfs := mock.ReqFactoryStatic(rss)
	rnfsrs := rnfs(ctx)

	assert.NotNil(t, rnfsrs, "Nil Req Static function was supposed to return a nil")
	assert.Equal(t, rss, rnfsrs, "Unexpected requirements returned")
}

func Test_MockDependencyDepStatic(t *testing.T) {
	ctx := context.Background()

	dnfn := mock.DepFactoryStatic(nil, nil)
	dsn, dsnerr := dnfn(ctx, nil)

	assert.Nil(t, dsnerr, "Nel Dep Static factory returned an unexpected error")
	assert.Nil(t, dsn, "Nil Dep Static factory was supposed to return a nil")

	d := mock.StaticDependency("this-dep", "expected dependency", nil, nil)
	dnfs := mock.DepFactoryStatic(d, nil)
	dnfsds, dnfsdserr := dnfs(ctx, nil)

	assert.Nil(t, dnfsdserr, "Dep Static factory returned an unexpected error")
	assert.Equal(t, d, dnfsds, "Unexpected dependency returned")
}

func Test_MockDependencyDepMatch(t *testing.T) {
	ctx := context.Background()

	pos := mock.SimpleRequirement("this", "expected match")
	neg := mock.SimpleRequirement("not-this", "expected to not match")

	d := mock.StaticDependency("this-dep", "expected dependency", nil, nil)
	dnfs := mock.DepFactoryIdMatch("this", d)

	dnfsposd, dnfsposderr := dnfs(ctx, pos)
	assert.Nil(t, dnfsposderr, "Dep match factory function produced unexpected error")
	assert.Equal(t, d, dnfsposd, "Expected dependency not returned for matching requirement")

	dnfsnegd, dnfsnegderr := dnfs(ctx, neg)
	assert.Nil(t, dnfsnegderr, "Dep match factory function produced unexpected error")
	assert.NotEqualValues(t, d, dnfsnegd, "Unexpected dependencies returned for not matching requirement")
}
