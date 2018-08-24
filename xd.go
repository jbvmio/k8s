package k8s

import (
	"encoding/json"
	"fmt"
	"time"
)

//XD contains the common, available data for all resources
type XD struct {
	*Metadata `json:"metadata,omitempty"`
	*Spec     `json:"spec,omitempty"`
	*Status   `json:"status,omitempty"`
}

// Metadata contains the available metadata for all resources
type Metadata struct {
	Name            string            `json:"name,omitempty"`
	Created         time.Time         `json:"creationTimestamp,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	Selflink        string            `json:"selfLink,omitempty"`
	UID             string            `json:"uid,omitempty"`
	OwnerReferences []OwnerReferences `json:"ownerReferences,omitempty"`
}

//OwnerReferences struct here.
type OwnerReferences struct {
	OwnerKind string `json:"kind,omitempty"`
	OwnerName string `json:"name,omitempty"`
	OwnerUID  string `json:"uid,omitempty"`
}

// Spec contains the available Spec data for all resources
type Spec struct {
	ClusterIP  string       `json:"clusterIP,omitempty"`
	ExternalID string       `json:"externalID,omitempty"`
	NodeName   string       `json:"nodeName,omitempty"`
	Replicas   int          `json:"replicas,omitempty"`
	Type       string       `json:"type,omitempty"`
	Containers []Containers `json:"containers,omitempty"`
	Ports      []Ports      `json:"ports,omitempty"`
}

// Containers struct definition:
type Containers struct {
	ContainerName  string  `json:"name,omitempty"`
	ContainerImage string  `json:"image,omitempty"`
	Ports          []Ports `json:"ports,omitempty"`
}

// Ports struct definition:
type Ports struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	Port          int    `json:"port,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	TargetPort    int    `json:"targetPort,omitempty"`
}

// Status contains the available Status data for all resources
type Status struct {
	HostIP             string              `json:"hostIP,omitempty"`
	Phase              string              `json:"phase,omitempty"`
	PodIP              string              `json:"podIP,omitempty"`
	AvailableReplicase int                 `json:"availableReplicas,omitempty"`
	ReadyReplicas      int                 `json:"readyReplicas,omitempty"`
	UpdatedReplicas    int                 `json:"updatedReplicas,omitempty"`
	ContainerStatuses  []ContainerStatuses `json:"containerStatuses,omitempty"`

	NodeInfo `json:"nodeInfo,omitempty"`
	//LoadBalancer       string `json:"loadBalancer,omitempty"`
}

// ContainerStatuses definition (is []array):
type ContainerStatuses struct {
	ContainerStatusName string `json:"name,omitempty"`
	Ready               bool   `json:"ready"`
	RestartCount        int    `json:"restartCount"`
	*State              `json:"state,omitempty"`
}

// State definition
type State struct {
	*Running    `json:"running,omitempty"`
	*Waiting    `json:"waiting,omitempty"`
	*Terminated `json:"terminated,omitempty"`
}

// Running definition
type Running struct {
	Started time.Time `json:"startedAt,omitempty"`
}

// Waiting definition
type Waiting struct {
	Message       string `json:"message,omitempty"`
	WaitingReason string `json:"reason,omitempty"`
}

// Terminated definition
type Terminated struct {
	Started          time.Time `json:"startedAt,omitempty"`
	Finished         time.Time `json:"finishedAt,omitempty"`
	TerminatedReason string    `json:"reason,omitempty"`
	ExitCode         int       `json:"exitCode,omitempty"`
}

// NodeInfo struct definition:
type NodeInfo struct {
	Architecture            string `json:"architecture,omitempty"`
	BootID                  string `json:"bootID,omitempty"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion,omitempty"`
	KernelVersion           string `json:"kernelVersion,omitempty"`
	KubeProxyVersion        string `json:"kubeProxyVersion,omitempty"`
	KubeletVersion          string `json:"kubeletVersion,omitempty"`
	MachineID               string `json:"machineID,omitempty"`
	OperatingSystem         string `json:"operatingSystem,omitempty"`
	OSImage                 string `json:"osImage,omitempty"`
	SystemUUID              string `json:"systemUUID,omitempty"`
}

// CreateXD takes json/output created by RawClient and creates an []XD struct if possible
func createXD(json []string) ([]XD, error) {
	var xds []XD
	var errd error
	xdChan := make(chan XD, 100)
	errChan := make(chan error, 100)
	for _, j := range json {
		go makeXD(j, xdChan, errChan)
	}
	for i := 0; i < len(json); i++ {
		select {
		case x := <-xdChan:
			xds = append(xds, x)
		case err := <-errChan:
			errd = err
		}
	}
	return xds, errd
}

func makeXD(data string, xdChan chan XD, errChan chan error) {
	xd := XD{}
	err := json.Unmarshal([]byte(data), &xd)
	if err != nil {
		errChan <- err
		return
	}
	xdChan <- xd
	return
}

// Raw returns the json output of XD
func (x *XD) Raw() string {
	j, _ := json.Marshal(x)
	return fmt.Sprintf("%s", j)
}
