package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KClient is the abstraction interface for interacting with k8s resources.
type KClient struct {
	cs       *kubernetes.Clientset
	ns       string
	listOpts metav1.ListOptions
}

// SetNS sets the namespace to operate within, defaults to all
func (kc *KClient) SetNS(ns string) {
	kc.ns = ns
}

// NewKClient returns a new RawClient
func NewKClient(inCluster bool) (*KClient, error) {
	var kc KClient
	if inCluster {
		cs, err := CreateICClientSet()
		if err != nil {
			return &kc, err
		}
		kc.cs = cs
		kc.listOpts = metav1.ListOptions{}
		return &kc, nil
	}
	cs, err := CreateOCClientSet()
	if err != nil {
		return &kc, err
	}
	kc.cs = cs
	kc.listOpts = metav1.ListOptions{}
	return &kc, nil
}
