package implementation_test

import (
	"context"
	"testing"

	mke3implementation "github.com/Mirantis/launchpad/mirantis/product/mke3/implementation"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

// Prove that the MKE3 Requirement and Dependency are valid and can be matched.
func Test_DependencySanity(t *testing.T) {
	ctx := t.Context()

	v := "3.7.7"
	apic := mke3implementation.Config{Version: v}
	dc := mke3implementation.MKE3DependencyConfig{Version: v}

	r := mke3implementation.NewMKE3Requirement(
		"test",
		"test mke3 requirement",
		dc,
	)

	var rr dependency.Requirement = r

	d := mke3implementation.NewMKE3Dependency(
		"test",
		"test mke3 dependency",
		func(context.Context) (*mke3implementation.API, error) {
			return mke3implementation.NewAPI(apic), nil
		},
	)

	var dd dependency.Dependency = d
	var _ mke3implementation.ProvidesMKE3 = d

	if _, ok := dd.(mke3implementation.ProvidesMKE3); !ok {
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
