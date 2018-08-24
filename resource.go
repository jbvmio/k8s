package k8s

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	listOpts = metav1.ListOptions{}
)

// Resource is the base representation for a resource located in a kubernetes cluster
type Resource struct {
	Name      string
	Namespace string
	Kind      string
	URL       string
}

// GetResources retreives the resource type as indicated by kind (ie. "pods") and returns an array of []Resource of the type.
func GetResources(kind string, cs *kubernetes.Clientset) ([]Resource, error) {
	var resources []Resource
	switch kind {
	case "pods":
		resource, err := cs.CoreV1().Pods("").List(listOpts)
		if err != nil {
			return resources, err
		}
		for _, i := range resource.Items {
			res := Resource{}
			res.Name = i.Name
			res.Namespace = i.Namespace
			res.Kind = "Pod"
			res.URL = i.SelfLink
			resources = append(resources, res)
		}
	case "nodes":
		resource, err := cs.CoreV1().Nodes().List(listOpts)
		if err != nil {
			return resources, err
		}
		for _, i := range resource.Items {
			res := Resource{}
			res.Name = i.Name
			res.Kind = "Node"
			res.URL = i.SelfLink
			resources = append(resources, res)
		}
	case "rs":
		resource, err := cs.ExtensionsV1beta1().ReplicaSets("").List(listOpts)
		if err != nil {
			return resources, err
		}
		for _, i := range resource.Items {
			res := Resource{}
			res.Name = i.Name
			res.Namespace = i.Namespace
			res.Kind = "ReplicaSet"
			res.URL = i.SelfLink
			resources = append(resources, res)
		}
	case "deploys":
		resource, err := cs.ExtensionsV1beta1().Deployments("").List(listOpts)
		if err != nil {
			return resources, err
		}
		for _, i := range resource.Items {
			res := Resource{}
			res.Name = i.Name
			res.Namespace = i.Namespace
			res.Kind = "Deployment"
			res.URL = i.SelfLink
			resources = append(resources, res)
		}
	default:
		return resources, fmt.Errorf("No resources found for the given request")
	}
	return resources, nil
}

// GetK8sValue takes a K8s item and retrieves a value for a given field (DevUse*)
func getK8sValue(field string, item interface{}) string {
	i := reflect.ValueOf(item)
	return i.FieldByName(field).String()
}
