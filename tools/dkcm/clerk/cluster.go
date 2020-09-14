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
)

// Cluster stores a row in the "Cluster" db table
type Cluster struct {
	*ClusterParams
	ProjectID string
	Status    string
}

// Function option that modify a field of Cluster
type ClusterOption func(*Cluster)

func NewCluster(opts ...ClusterOption) *Cluster {
	c := &Cluster{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c Cluster) String() string {
	return fmt.Sprintf("Cluster Info: (ProjectID: %s, NodesCount: %d, NodeType: %s, Status: %s, Zone: %s)",
		c.ProjectID, c.Nodes, c.NodeType, c.Status, c.Zone)
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
