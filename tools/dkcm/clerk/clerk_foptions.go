/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package clerk

import (
	"fmt"
	"time"
)

// Function option that returns a part of query for elements in ClusterParams
type QueryClusterParamsOption func(*ClusterParams) string

// Function option that returns as a key value pair for database update
type UpdateOption func() string

// Function option that modify a field of Request
type RequestOption func(*Request)

// Function option that modify a field of Cluster
type ClusterOption func(*Cluster)

// Function option that modify a field of ClusterParams
type ClusterParamsOption func(*ClusterParams)

func UpdateStringField(key string, value string) UpdateOption {
	return func() string {
		return fmt.Sprintf("%s = %s", key, value)
	}
}

func UpdateNumField(key string, value int64) UpdateOption {
	return func() string {
		return fmt.Sprintf("%s = %v", key, value)
	}
}

func QueryZone() QueryClusterParamsOption {
	return func(cp *ClusterParams) string {
		return fmt.Sprintf("Zone = %s", cp.Zone)
	}
}

func QueryNodes() QueryClusterParamsOption {
	return func(cp *ClusterParams) string {
		return fmt.Sprintf("Nodes = %d", cp.Nodes)
	}
}

func QueryNodeType() QueryClusterParamsOption {
	return func(cp *ClusterParams) string {
		return fmt.Sprintf("NodeType = %s", cp.NodeType)
	}
}

func AddZone(zone string) ClusterParamsOption {
	return func(cp *ClusterParams) {
		cp.Zone = zone
	}
}

func AddNodes(nodes int64) ClusterParamsOption {
	return func(cp *ClusterParams) {
		cp.Nodes = nodes
	}
}

func AddNodeType(nodeType string) ClusterParamsOption {
	return func(cp *ClusterParams) {
		cp.NodeType = nodeType
	}
}

func AddClusterID(clusterID int64) ClusterOption {
	return func(c *Cluster) {
		c.ClusterID = clusterID
	}
}

func AddProjectID(projectID string) ClusterOption {
	return func(c *Cluster) {
		c.ProjectID = projectID
	}
}

func AddStatus(status string) ClusterOption {
	return func(c *Cluster) {
		c.Status = status
	}
}

func AddaccessToken(accessToken string) RequestOption {
	return func(r *Request) {
		r.accessToken = accessToken
	}
}

func AddrequestTime(requestTime time.Time) RequestOption {
	return func(r *Request) {
		r.requestTime = requestTime
	}
}

func AddProwJobID(prowJobID string) RequestOption {
	return func(r *Request) {
		r.ProwJobID = prowJobID
	}
}

func AddRClusterID(clusterID int64) RequestOption {
	return func(r *Request) {
		r.ClusterID = clusterID
	}
}
