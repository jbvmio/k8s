// url.go is WiP for experimental and testing purposes only.

package k8s

import (
	"fmt"
	"strings"
)

// URL is the selfLink target for a k8s resource
type URL struct {
	SelfLink string
}

// SelfLinker implements the selfLink return interface
type SelfLinker interface {
	FindSelfLinks(kind, name string) ([]string, error)
}

// URLGetter implements the []URL return interface
type URLGetter interface {
	GetAllURLs(s string) ([]URL, error)
}

// GetAllURLs searches for URLs using either a k8s KClient or RawClient
func GetAllURLs(client URLGetter, s string) ([]URL, error) {
	urls, err := client.GetAllURLs(s)
	return urls, err
}

// FindSelfLinks returns a []string of selfLinks
func FindSelfLinks(client SelfLinker, kind, name string) ([]string, error) {
	urls, err := client.FindSelfLinks(kind, name)
	return urls, err
}

// ToResource converts a URL into a Resource (ToDo*)
func (url *URL) ToResource() {

}

// FindSelfLinks returns all resources matching the kind and name specified
func (rc *RawClient) FindSelfLinks(kind, name string) ([]string, error) {
	var results []string
	all, err := rc.RawAll(kind)
	if err != nil {
		return results, err
	}
	results = findSelfLinks(all, name)
	return results, nil
}

// GetAllURLs retreives all selfLink URLs for the specified kind as indicated by k (ie. "pods")
func (rc *RawClient) GetAllURLs(kind string) ([]URL, error) {
	var urls []URL
	all, err := rc.RawAll(kind)
	if err != nil {
		return urls, err
	}
	links := returnSelfLinks(all)
	for _, l := range links {
		url := URL{
			SelfLink: l,
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// GetAllURLs retreives all selfLink URLs for the specified kind as indicated by kind (ie. "pods")
func (kc *KClient) GetAllURLs(kind string) ([]URL, error) {
	var urls []URL
	switch kind {
	case "pods":
		resource, err := kc.cs.CoreV1().Pods(kc.ns).List(kc.listOpts)
		if err != nil {
			return urls, err
		}
		for _, i := range resource.Items {
			url := URL{}
			url.SelfLink = i.SelfLink
			urls = append(urls, url)
		}
	case "nodes":
		resource, err := kc.cs.CoreV1().Nodes().List(kc.listOpts)
		if err != nil {
			return urls, err
		}
		for _, i := range resource.Items {
			url := URL{}
			url.SelfLink = i.SelfLink
			urls = append(urls, url)
		}
	case "rs":
		resource, err := kc.cs.ExtensionsV1beta1().ReplicaSets(kc.ns).List(kc.listOpts)
		if err != nil {
			return urls, err
		}
		for _, i := range resource.Items {
			url := URL{}
			url.SelfLink = i.SelfLink
			urls = append(urls, url)
		}
	case "deploys":
		resource, err := kc.cs.ExtensionsV1beta1().Deployments(kc.ns).List(kc.listOpts)
		if err != nil {
			return urls, err
		}
		for _, i := range resource.Items {
			url := URL{}
			url.SelfLink = i.SelfLink
			urls = append(urls, url)
		}
	default:
		return urls, fmt.Errorf("No resources found for the given request")
	}
	return urls, nil
}

// FindSelfLinks returns all resources matching the kind and name specified
func (kc *KClient) FindSelfLinks(kind, name string) ([]string, error) {
	var results []string
	switch kind {
	case "pods":
		resource, err := kc.cs.CoreV1().Pods(kc.ns).List(kc.listOpts)
		if err != nil {
			return results, err
		}
		for _, i := range resource.Items {
			if strings.Contains(i.Name, name) {
				results = append(results, i.SelfLink)
			}
		}
	case "nodes":
		resource, err := kc.cs.CoreV1().Nodes().List(kc.listOpts)
		if err != nil {
			return results, err
		}
		for _, i := range resource.Items {
			if strings.Contains(i.Name, name) {
				results = append(results, i.SelfLink)
			}
		}
	default:
		return results, fmt.Errorf("No resources found for the given request")
	}
	return results, nil
}
