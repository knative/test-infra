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

// Function option that returns a part of query for elements in ClusterParams
type QueryClusterParamsOption func(*ClusterParams) string

// Function option that returns as a key value pair for database update
type UpdateOption func() string

func UpdateStringField(key string, value string) UpdateOption {
	return func() string {
		return fmt.Sprintf("%s = '%s'", key, value)
	}
}

func UpdateNumField(key string, value int64) UpdateOption {
	return func() string {
		return fmt.Sprintf("%s = %d", key, value)
	}
}

func QueryZone() QueryClusterParamsOption {
	return func(cp *ClusterParams) string {
		return fmt.Sprintf("Zone = '%s'", cp.Zone)
	}
}

func QueryNodes() QueryClusterParamsOption {
	return func(cp *ClusterParams) string {
		return fmt.Sprintf("Nodes = %d", cp.Nodes)
	}
}

func QueryNodeType() QueryClusterParamsOption {
	return func(cp *ClusterParams) string {
		return fmt.Sprintf("NodeType = '%s'", cp.NodeType)
	}
}
