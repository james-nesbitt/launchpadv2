package kubernetes_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

func Test_DependencySanity(t *testing.T) {
	ctx := context.Background()
	cfg := kubernetes.Config{}
	k := kubernetes.NewKubernetes(cfg)
	v := kubernetes.Version{Version: "1.25"}
	kr := kubernetes.NewKubernetesRequirement("test", "my test kubernetes requirement", v)
	kd := kubernetes.NewKubernetesDependency("test", "my test kubernetes dependency", func(context.Context) (*kubernetes.Kubernetes, error) {
		return k, nil
	})

	var r dependency.Requirement = kr

	kr2, ok := r.(kubernetes.KubernetesRequirement)
	if !ok {
		t.Errorf("kubernetes requirement is not a valid kubernetes requirement: %+v", r)
	}

	v2 := kr2.RequiresKubernetes(ctx)
	if v.Version != v2.Version {
		t.Errorf("Wrong version given: %+v", v2)
	}

	if err := r.Match(kd); err != nil {
		t.Errorf("Assigning kube dep to kube req returned an error: %s", err.Error())
	}

	d := r.Matched(ctx)
	if d == nil {
		t.Fatalf("Kubernetes requirement didn't return dependency after being matched: %+v", r)
	}

	kd2, ok := d.(kubernetes.KubernetesDependency)
	if !ok {
		t.Fatalf("Kubernetes requirement returned unexpected dependency: %+v", kd2)
	}

	k2 := kd2.Kubernetes(ctx)

	if k2 != k {
		t.Errorf("Kubernetes dependency returned unexpected kubernetes client: %+v", k2)
	}

	if d.Id() != "test" {
		t.Errorf("Kubernetes dependency returned wrong id: %s", d.Id())
	}
}
