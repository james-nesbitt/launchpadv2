package mke3_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/product/mke3"
)

func Test_DependencySanity(t *testing.T) {
	ctx := context.Background()

	mdc := mke3.MKE3DependencyConfig{
		Version: "3.7.7",
	}
	m := mke3.NewMKE3("test", mke3.Config{Version: mdc.Version})

	r := mke3.NewMKE3Requirement(
		"test",
		"test mke3 requirement",
		mdc,
	)

	var rr dependency.Requirement = r

	d := mke3.NewMKE3Dependency(
		"test",
		"test mke3 dependency",
		func(context.Context) (*mke3.MKE3, error) {
			return m, nil
		},
	)

	var dd dependency.Dependency = d
	var _ mke3.MKE3Dependency = d

	if _, ok := dd.(mke3.MKE3Dependency); !ok {
		t.Errorf("Could not convert our dependency to the MKE3 Dependency: %+v", dd)
	}

	if db := rr.Matched(ctx); db != nil {
		t.Errorf("requirements says it is matched before matching: %+v", r)
	}

	if err := rr.Match(dd); err != nil {
		t.Errorf("requirement failed matching with our dependency: %s \n%+v \n %+v", err.Error(), r, d)
	}

	dd2 := rr.Matched(ctx)
	if dd2 == nil {
		t.Errorf("requirements says it isn't matched after matching: %+v", r)
	}
}
