package kubernetes

import (
	"context"
	"fmt"

	kubeapi_errors "k8s.io/apimachinery/pkg/api/errors"
	kube_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeapi_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubeapi_schema "k8s.io/apimachinery/pkg/runtime/schema"
	kube_dynamic "k8s.io/client-go/dynamic"
	kube_kubernetes "k8s.io/client-go/kubernetes"
)

func (k Kubernetes) ApplyUnstructuredObject(ctx context.Context, obj *kubeapi_unstructured.Unstructured) error {
	cs, cserr := k.Client()
	if cserr != nil {
		return cserr
	}
	dc, dcerr := k.DynamicClient()
	if dcerr != nil {
		return dcerr
	}

	return createOrUpdateObject(ctx, cs, dc, obj)
}

func (k Kubernetes) DeleteUnstructuredObject(ctx context.Context, obj *kubeapi_unstructured.Unstructured) error {
	cs, cserr := k.Client()
	if cserr != nil {
		return cserr
	}
	dc, dcerr := k.DynamicClient()
	if dcerr != nil {
		return dcerr
	}

	return deleteObject(ctx, cs, dc, obj)
}

/** Unstructured handling

stolen from:  https://github.com/MirantisContainers/blueprint-cli/blob/main/pkg/k8s/dynamic.go

*/

func createOrUpdateObject(ctx context.Context, client kube_kubernetes.Interface, dynamicClient kube_dynamic.Interface, obj *kubeapi_unstructured.Unstructured) error {
	gvr, _ := getResource(client, obj)
	namespace := obj.GetNamespace()
	objName := obj.GetName()

	existing, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, objName, kube_metav1.GetOptions{})
	if kubeapi_errors.IsNotFound(err) {
		_, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, obj, kube_metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create resource: %q", err)
		}
	} else {
		obj.SetResourceVersion(existing.GetResourceVersion())
		_, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(ctx, obj, kube_metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update resource: %q", err)
		}
	}

	return nil
}

func deleteObject(ctx context.Context, client kube_kubernetes.Interface, dynamicClient kube_dynamic.Interface, obj *kubeapi_unstructured.Unstructured) error {
	gvr, _ := getResource(client, obj)
	namespace := obj.GetNamespace()
	objName := obj.GetName()

	_, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, objName, kube_metav1.GetOptions{})
	if kubeapi_errors.IsNotFound(err) {
		return nil
	}

	err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, obj.GetName(), kube_metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete resource: %q", err)
	}

	return nil
}

// getResource returns the GroupVersionResource for a given object
// Especially, it discovers the resource name for the given object
func getResource(client kube_kubernetes.Interface, object *kubeapi_unstructured.Unstructured) (kubeapi_schema.GroupVersionResource, error) {
	gvk := object.GroupVersionKind()
	apiResourceList, err := client.Discovery().ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return kubeapi_schema.GroupVersionResource{}, err
	}

	var resource *kube_metav1.APIResource
	for _, r := range apiResourceList.APIResources {
		if r.Kind == gvk.Kind {
			resource = &r
			break
		}
	}

	return kubeapi_schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: resource.Name}, nil
}
