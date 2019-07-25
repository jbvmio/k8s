package k8s

import (
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// RawClient is like Client but dedicated for raw input/output
type RawClient struct {
	cs           *kubernetes.Clientset
	ns           string
	exactMatches bool
	listOpts     metav1.ListOptions
	delOpts      metav1.DeleteOptions
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

// NewManualClient returns a new RawClient using just a the API url
func NewManualClient(url string) (*RawClient, error) {
	var rc RawClient
	cs, err := CreateManualClientSet(url)
	if err != nil {
		return &rc, err
	}
	rc.cs = cs
	rc.listOpts = metav1.ListOptions{}
	return &rc, nil
}

// ClientFromConfig returns a new RawClient using a manual config
func ClientFromConfig(config *rest.Config) (*RawClient, error) {
	var rc RawClient
	config.ContentConfig.GroupVersion = &v1.SchemeGroupVersion
	config.ContentConfig.NegotiatedSerializer = scheme.Codecs
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &rc, err
	}
	rc.cs = cs
	rc.listOpts = metav1.ListOptions{}
	return &rc, nil
}
