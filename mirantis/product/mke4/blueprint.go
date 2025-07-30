package mke4

import (
	"context"
	"errors"
	"fmt"

	kubeapi_errors "k8s.io/apimachinery/pkg/api/errors"
	kubeapi_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubernetes "github.com/Mirantis/launchpad/implementation/kubernetes"
)

// @SEE https://github.com/MirantisContainers/blueprint-cli/blob/main/pkg/constants/constants.go#L5
// @TODO consider importing it instead of duplicating it (if they OSS it).
const (
	// ManifestUrlLatest is the URL of the latest manifest YAML for the Boundless Operator.
	ManifestUrlLatest = "https://raw.githubusercontent.com/mirantiscontainers/boundless/main/deploy/static/boundless-operator.yaml"

	// NamespaceBlueprint is the system namespace where the Boundless Operator and its components are installed.
	NamespaceBlueprint = "blueprint-system"

	// DefaultBlueprintFileName represents the default blueprint filename.
	DefaultBlueprintFileName = "blueprint.yaml"

	// BlueprintOperatorDeployment if thi sdeployment name is running, then MKE4 is active.
	BlueprintOperatorDeployment = "blueprint-operator-controller-manager"
)

var (
	ErrOperatorInstall    = errors.New("could not install MKE4 operator")
	ErrOperatorUninstall  = errors.New("could not uninstall MKE4 operator")
	ErrBlueprintInstall   = errors.New("could not install MKE4 instance")
	ErrBlueprintUninstall = errors.New("could not uninstall MKE4 instance")
)

// OperaterInstall install the MKE4 operator onto the cluster.
func OperatorInstall(ctx context.Context, k kubernetes.Kubernetes, c Config) error {
	ors, orserr := c.OperatorResources()
	if orserr != nil {
		return fmt.Errorf("%w, %s", ErrOperatorInstall, orserr)
	}

	for _, or := range ors {
		if err := k.ApplyUnstructuredObject(ctx, &or); err != nil {
			return fmt.Errorf("%w, %s", ErrOperatorInstall, err)
		}
	}

	return nil
}

// OperatorUninstall uninstall the MKE4 operator from the cluster.
func OperatorUninstall(ctx context.Context, k kubernetes.Kubernetes, c Config) error {
	ors, orserr := c.OperatorResources()
	if orserr != nil {
		return fmt.Errorf("%w, %s", ErrOperatorInstall, orserr)
	}

	for _, or := range ors {
		if err := k.DeleteUnstructuredObject(ctx, &or); err != nil {
			return fmt.Errorf("%w, %s", ErrOperatorUninstall, err)
		}
	}

	return nil
}

// ActivateOrUpdate creates or updates a kubernetes CR for MKE4.
func ActivateOrUpdate(ctx context.Context, k kubernetes.Kubernetes, c Config) error {
	bp, bperr := c.BlueprintResource()
	if bperr != nil {
		return fmt.Errorf("%w, %s", ErrBlueprintInstall, bperr)
	}

	if err := k.ApplyUnstructuredObject(ctx, &bp); err != nil {
		return fmt.Errorf("%w, %s", ErrBlueprintInstall, err)
	}

	return nil
}

// Deactivate deletes the CR instance for MKE4.
func Deactivate(ctx context.Context, k kubernetes.Kubernetes, c Config) error {
	bp, bperr := c.BlueprintResource()
	if bperr != nil {
		return fmt.Errorf("%w, %s", ErrBlueprintUninstall, bperr)
	}

	if err := k.DeleteUnstructuredObject(ctx, &bp); err != nil {
		return fmt.Errorf("%w, %s", ErrBlueprintUninstall, err)
	}

	return nil
}

// IsActive is the MKE4 operator active.
func IsActive(ctx context.Context, k kubernetes.Kubernetes, c Config) (bool, error) {
	kc, kcerr := k.Client()
	if kcerr != nil {
		return false, kcerr
	}

	_, derr := kc.AppsV1().Deployments(c.Namespace).Get(ctx, BlueprintOperatorDeployment, kubeapi_metav1.GetOptions{})
	if derr != nil {
		if kubeapi_errors.IsNotFound(derr) {
			return false, nil
		} else {
			return false, derr
		}
	}
	return true, nil
}
