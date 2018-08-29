package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// RawClient is like Client but dedicated for raw input/output
type RawClient struct {
	cs           *kubernetes.Clientset
	ns           string
	exactMatches bool
	listOpts     metav1.ListOptions
}

// SetNS sets the namespace to operate within, defaults to all
func (rc *RawClient) SetNS(ns string) {
	rc.ns = ns
}

// ExactMatches sets the behavior of searching for resources, either wildcard or exact, defaults to wildcard or false
func (rc *RawClient) ExactMatches(x bool) {
	rc.exactMatches = x
}

// NewRawClient returns a new RawClient
func NewRawClient(inCluster bool) (*RawClient, error) {
	var rc RawClient
	if inCluster {
		cs, err := CreateICClientSet()
		if err != nil {
			return &rc, err
		}
		rc.cs = cs
		rc.listOpts = metav1.ListOptions{}
		return &rc, nil
	}
	cs, err := CreateOCClientSet()
	if err != nil {
		return &rc, err
	}
	rc.cs = cs
	rc.listOpts = metav1.ListOptions{}
	return &rc, nil
}

// RawAllPods returns all pods in json
func (rc *RawClient) RawAllPods() ([]byte, error) {
	p, err := rc.cs.CoreV1().RESTClient().Get().Namespace(rc.ns).Resource("pods").VersionedParams(&rc.listOpts, scheme.ParameterCodec).DoRaw()
	return p, err
}

// RawAllNodes returns all nodes in json
func (rc *RawClient) RawAllNodes() ([]byte, error) {
	n, err := rc.cs.CoreV1().RESTClient().Get().Resource("nodes").VersionedParams(&rc.listOpts, scheme.ParameterCodec).DoRaw()
	return n, err
}

// RawAllSvc returns all services in json
func (rc *RawClient) RawAllSvc() ([]byte, error) {
	d, err := rc.cs.CoreV1().RESTClient().Get().Namespace(rc.ns).Resource("services").VersionedParams(&rc.listOpts, scheme.ParameterCodec).DoRaw()
	return d, err
}

// RawAllDeployments returns all deployments in json
func (rc *RawClient) RawAllDeployments() ([]byte, error) {
	d, err := rc.cs.ExtensionsV1beta1().RESTClient().Get().Namespace(rc.ns).Resource("deployments").VersionedParams(&rc.listOpts, scheme.ParameterCodec).DoRaw()
	return d, err
}

// RawAllRS returns all replicasets in json
func (rc *RawClient) RawAllRS() ([]byte, error) {
	d, err := rc.cs.ExtensionsV1beta1().RESTClient().Get().Namespace(rc.ns).Resource("replicasets").VersionedParams(&rc.listOpts, scheme.ParameterCodec).DoRaw()
	return d, err
}

// RawAllIngress returns all ingresses in json
func (rc *RawClient) RawAllIngress() ([]byte, error) {
	d, err := rc.cs.ExtensionsV1beta1().RESTClient().Get().Namespace(rc.ns).Resource("ingresses").VersionedParams(&rc.listOpts, scheme.ParameterCodec).DoRaw()
	return d, err
}

// RawAll returns all of a given kind in json
func (rc *RawClient) RawAll(kind string) ([]byte, error) {
	var all []byte
	var errd error
	switch kind {
	case "pods":
		all, errd = rc.RawAllPods()
		if errd != nil {
			return all, errd
		}
	case "nodes":
		all, errd = rc.RawAllNodes()
		if errd != nil {
			return all, errd
		}
	case "replicasets":
		all, errd = rc.RawAllRS()
		if errd != nil {
			return all, errd
		}
	case "deployments":
		all, errd = rc.RawAllDeployments()
		if errd != nil {
			return all, errd
		}
	case "services":
		all, errd = rc.RawAllSvc()
		if errd != nil {
			return all, errd
		}
	case "ingresses":
		all, errd = rc.RawAllIngress()
		if errd != nil {
			return all, errd
		}
	default:
		return all, fmt.Errorf("No resources found for the given request")
	}
	return all, nil
}

// RawNodes is a convenience wrapper for SearchFor returning all nodes matching the specified name
func (rc *RawClient) RawNodes(name string) ([]string, error) {
	results, err := rc.SearchFor("nodes", name)
	return results, err
}

// RawPods is a convenience wrapper for SearchFor returning all pods matching the specified name
func (rc *RawClient) RawPods(name string) ([]string, error) {
	results, err := rc.SearchFor("pods", name)
	return results, err
}

// RawRS is a convenience wrapper for SearchFor returning all replicasets matching the specified name
func (rc *RawClient) RawRS(name string) ([]string, error) {
	results, err := rc.SearchFor("replicasets", name)
	return results, err
}

// RawDeployments is a convenience wrapper for SearchFor returning all deployments matching the specified name
func (rc *RawClient) RawDeployments(name string) ([]string, error) {
	results, err := rc.SearchFor("deployments", name)
	return results, err
}

// RawSvc is a convenience wrapper for SearchFor returning all services matching the specified name
func (rc *RawClient) RawSvc(name string) ([]string, error) {
	results, err := rc.SearchFor("services", name)
	return results, err
}

// RawIngress is a convenience wrapper for SearchFor returning all ingresses matching the specified name
func (rc *RawClient) RawIngress(name string) ([]string, error) {
	results, err := rc.SearchFor("ingresses", name)
	return results, err
}

// SearchFor returns resources matching the kind and name specified as an array containing json of each result
func (rc *RawClient) SearchFor(kind, name string) ([]string, error) {
	var jsonArray []string
	all, err := rc.RawAll(kind)
	if err != nil {
		return jsonArray, err
	}
	if rc.exactMatches {
		jsonArray = parseExact(all, name)
		return jsonArray, nil
	}
	jsonArray = parseFor(all, name)
	return jsonArray, nil
}
