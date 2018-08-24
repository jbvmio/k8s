package k8s

import (
	"fmt"

	"k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// ConfigLocation sets the location for kubeconfig
	ConfigLocation        string
	defaultConfigLocation = string(homeDir() + "/.kube/config")
)

// ManualConfig starts and returns a rest.Config using a master url
// Any other options and authentication will need to be set manually
func ManualConfig(url string) (*rest.Config, error) {
	var config *rest.Config
	config, err := clientcmd.BuildConfigFromFlags(url, "")
	if err != nil {
		return config, err
	}
	return config, nil
}

// CreateClientSet Creates a Clientset from a rest.Config
func CreateClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, err
	}
	return clientset, nil
}

// GetKubeConfig returns a rest.Config using the current-context set in kubeconfig
// Use ConfigLocation variable to set the kubeconfig location if not located in $HOME/.kube/config
func GetKubeConfig() (*rest.Config, error) {
	var config *rest.Config
	var kubeconfig string
	if ConfigLocation == "" {
		if fileExists(defaultConfigLocation) {
			kubeconfig = defaultConfigLocation
		} else {
			return config, fmt.Errorf("cannot locate kubeconfig at %v", defaultConfigLocation)
		}
	} else {
		kubeconfig = ConfigLocation
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return config, err
	}
	return config, nil
}

// CreateRestClient returns a new RestClient (WiP*)
func createRestClient() (*rest.RESTClient, error) {
	var restClient *rest.RESTClient
	config, err := GetKubeConfig()
	if err != nil {
		return restClient, err
	}

	config.ContentConfig.GroupVersion = &v1.SchemeGroupVersion
	config.ContentConfig.NegotiatedSerializer = scheme.Codecs

	restClient, err = rest.RESTClientFor(config)
	if err != nil {
		return restClient, err
	}
	return restClient, nil
}

// CreateOCClientSet returns an Out of Cluster Clientset
func CreateOCClientSet() (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	config, err := GetKubeConfig()
	if err != nil {
		return clientset, err
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, err
	}
	return clientset, nil
}

// CreateICClientSet returns an In Cluster Clientset
func CreateICClientSet() (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	config, err := rest.InClusterConfig()
	if err != nil {
		return clientset, err
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, err
	}
	return clientset, nil
}
