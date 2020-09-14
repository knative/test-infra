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

// ClusterParams is a struct for common cluster parameters shared by both Cluster and Request
type ClusterParams struct {
	ID       int64
	Zone     string
	Nodes    int64
	NodeType string
}

// Function option that modify a field of ClusterParams
type ClusterParamsOption func(*ClusterParams)

func NewClusterParams(opts ...ClusterParamsOption) *ClusterParams {
	cp := &ClusterParams{}
	for _, opt := range opts {
		opt(cp)
	}
	return cp
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
