package mke3_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/product/mke3"
	mke3implementation "github.com/Mirantis/launchpad/pkg/product/mke3/implementation"
)

// Prove that the MKE3 component can handle MKE3 requirements
//
//	This is the mechanism by which requirements and dependencies are matched in MKE3
func Test_ProvidesMKE3(t *testing.T) {
	ctx := context.Background()

	dc := mke3implementation.MKE3DependencyConfig{Version: "3.7.7"}
	m := mke3.NewComponent("test", mke3.Config{Version: dc.Version})
	var mc component.Component = m

	r := mke3implementation.NewMKE3Requirement(
		"test",
		"test mke3 requirement",
		dc,
	)

	mp, ok := mc.(dependency.ProvidesDependencies)
	if !ok {
		t.Fatal("MKE3 component does not provide dependencies")
	}

	d, err := mp.ProvidesDependencies(ctx, r)
	if err != nil {
		t.Fatalf("MKE3 could not provide a dependency for an MKE requirement: %s", err.Error())
	}

	dm, ok := d.(mke3implementation.ProvidesMKE3)
	if !ok {
		t.Fatalf("MKE3 returned an invalid dependency: %+v", d)
	}

	api := dm.ProvidesMKE3(ctx)
	if api == nil {
		t.Errorf("MKE3 dependency did not return an MKE3 instance")
	}
}

// Prove that the MKE3 component can handle K8s requirements.
func Test_ProvidesK8s(t *testing.T) {
	ctx := context.Background()

	m := mke3.NewComponent("test", mke3.Config{Version: "3.7.7"})
	var mc component.Component = m

	kver := kubernetes.Version{Version: "1.28"}
	r := kubernetes.NewKubernetesRequirement(
		"test",
		"something needs kubernetes, MKE3 should provide that",
		kver,
	)

	mp, ok := mc.(dependency.ProvidesDependencies)
	if !ok {
		t.Fatal("MKE3 component does not provide dependencies")
	}

	d, err := mp.ProvidesDependencies(ctx, r)
	if err != nil {
		t.Fatalf("MKE3 could not provide a k8s dependency for a K8s requirement: %s", err.Error())
	}

	dm, ok := d.(kubernetes.KubernetesDependency)
	if !ok {
		t.Fatalf("MKE3 returned an invalid k8s dependency: %+v", d)
	}

	k8s := dm.Kubernetes(ctx)
	if k8s == nil {
		t.Errorf("empty k8s returned by kubernetes dependency")
	}
}
