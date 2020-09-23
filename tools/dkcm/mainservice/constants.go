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

package mainservice

const (
	// default parameters for mainservice server
	DefaultClusterName   = "e2e-cls"
	DefaultNetworkName   = "e2e-network"
	DefaultZone          = "us-west1"
	DefaultNodeType      = "e2-standard-4"
	DefaultOverProvision = 5
	DefaultNodesCount    = 4
	DefaultTimeOut       = 60
	DefaultPort          = "8080"

	// four statuses of cluster
	Ready = "Ready"
	WIP   = "WIP"
	InUse = "In Use"
	Fail  = "Failed"

	// Field for use when query Cluster db
	Status    = "Status"
	ClusterID = "ClusterID"

	// time interval to examine timeout requests
	CheckInterval = 2
)
