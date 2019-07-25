package k8s

// Delete removes a Resource
/*
func (rc *RawClient) Delete(r Resource) (Results, error) {
	res, err := rc.findManyResults("ingresses", search)
	if err != nil {
		return res, err
	}
	return res, nil
}
*/

// RawDeleteRS removes the named ReplicaSet
func (rc *RawClient) RawDeleteRS(name string) ([]byte, error) {
	d, err := rc.cs.ExtensionsV1beta1().RESTClient().Delete().Namespace(rc.ns).Resource("replicasets").Name(name).Body(&rc.delOpts).DoRaw()
	return d, err
}

/*
// RawDelete removes the specified kind by the name provided
func (rc *RawClient) rawDelete(kind, name string) error {
	switch kind {
	case "replicasets":
		err := rc.cs.ExtensionsV1beta1().RESTClient().Delete().Namespace(rc.ns).Resource("replicasets").Name("blahblah").Body(&rc.delOpts).DoRaw()
		return err

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
		return fmt.Errorf("No resources found for the given request")
	}

}
*/
