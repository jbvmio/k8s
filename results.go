package k8s

import "fmt"

// Results is the base construction for a search result
type Results struct {
	Kind  string
	XData []XD
}

// GetContainers returns a Containers Result by converting a Pod Result (WiP*)
func (r *Results) GetContainers() (Results, error) {
	var results Results
	if r.Kind != "Pod" {
		return results, fmt.Errorf("cannot convert Results of kind %v into a Containers", r.Kind)
	}
	return results, nil
}

// GetMany returns Results from the given kind and a list of names contained in an string array
// kind options are "pods", "nodes", "replicasets", "deployments", "services", "ingresses"
func (rc *RawClient) GetMany(kind string, names []string) (Results, error) {
	res, err := rc.findManyResults(kind, names)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetPods return pod based Results based on the search strings
func (rc *RawClient) GetPods(search ...string) (Results, error) {
	//res, err := rc.findResults("pods", search)
	res, err := rc.findManyResults("pods", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetNodes returns node based Results based on the search string
func (rc *RawClient) GetNodes(search string) (Results, error) {
	//res, err := rc.findResults("nodes", search)
	res, err := rc.findManyResults("nodes", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetRS returns replicaset based Results based on the search string
func (rc *RawClient) GetRS(search string) (Results, error) {
	//res, err := rc.findResults("replicasets", search)
	res, err := rc.findManyResults("replicasets", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetSvc returns node service Results based on the search string
func (rc *RawClient) GetSvc(search string) (Results, error) {
	//res, err := rc.findResults("services", search)
	res, err := rc.findManyResults("services", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetDeployments returns deployment based Results based on the search string
func (rc *RawClient) GetDeployments(search string) (Results, error) {
	//res, err := rc.findResults("deployments", search)
	res, err := rc.findManyResults("deployments", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetIngress returns ingress based Results based on the search string
func (rc *RawClient) GetIngress(search string) (Results, error) {
	//res, err := rc.findResults("ingresses", search)
	res, err := rc.findManyResults("ingresses", search)
	if err != nil {
		return res, err
	}
	return res, nil
}

// findResults returns Results based on the kind and target name search string
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

// findManyResults returns Results based on the kind and target names contained in an array []string
func (rc *RawClient) findManyResults(kind string, names []string) (Results, error) {
	var res Results
	var errd error
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
	var data []XD
	xdChan := make(chan []XD, 100)
	errChan := make(chan error, 10)
	all, err := rc.RawAll(kind)
	if err != nil {
		return res, err
	}
	for _, n := range names {
		go rc.getXD(n, &all, xdChan, errChan)
	}
	for i := 0; i < len(names); i++ {
		select {
		case err := <-errChan:
			errd = err
		case xd := <-xdChan:
			for _, x := range xd {
				data = append(data, x)
			}
		}
	}
	res.XData = data
	return res, errd
}

func (rc *RawClient) getXD(name string, all *[]byte, xdChan chan []XD, errChan chan error) {
	found := parseFor(*all, name)
	data, err := createXD(found)
	if err != nil {
		errChan <- err
		return
	}
	xdChan <- data
	return
}
