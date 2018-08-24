package k8s

import "fmt"

// Results is the base construction for a search result
type Results struct {
	Kind  string
	XData []XD
}

// GetPods returns pod based Results based on the search string
func (rc *RawClient) GetPods(search string) (Results, error) {
	res, err := rc.findResults("pods", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetNodes returns node based Results based on the search string
func (rc *RawClient) GetNodes(search string) (Results, error) {
	res, err := rc.findResults("nodes", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetRS returns replicaset based Results based on the search string
func (rc *RawClient) GetRS(search string) (Results, error) {
	res, err := rc.findResults("replicasets", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetSvc returns node service Results based on the search string
func (rc *RawClient) GetSvc(search string) (Results, error) {
	res, err := rc.findResults("services", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetDeployments returns deployment based Results based on the search string
func (rc *RawClient) GetDeployments(search string) (Results, error) {
	res, err := rc.findResults("deployments", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetIngress returns ingress based Results based on the search string
func (rc *RawClient) GetIngress(search string) (Results, error) {
	res, err := rc.findResults("ingresses", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// findResults returns pod based Results based on the search string
func (rc *RawClient) findResults(kind, search string) (Results, error) {
	var res Results
	switch kind {
	case "pods":
		res.Kind = "Pod"
	case "nodes":
		res.Kind = "Node"
	case "replicasets":
		res.Kind = "ReplicaSet"
	case "deployments":
		res.Kind = "Deployment"
	case "services":
		res.Kind = "Service"
	case "ingresses":
		res.Kind = "Ingress"
	default:
		return res, fmt.Errorf("No resources found for the given request")
	}
	found, err := rc.SearchFor(kind, search)
	if err != nil {
		return res, err
	}
	data, err := createXD(found)
	if err != nil {
		return res, err
	}
	res.XData = data
	return res, nil
}
