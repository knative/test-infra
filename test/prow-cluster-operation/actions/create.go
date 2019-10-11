/*
Copyright 2019 The Knative Authors

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

package actions

import (
	"log"
	"strconv"
	"strings"

	container "google.golang.org/api/container/v1beta1"
	"knative.dev/pkg/testutils/clustermanager"
	"knative.dev/pkg/testutils/common"
	"knative.dev/test-infra/test/metahelper/util"
	"knative.dev/test-infra/test/prow-cluster-operation/options"
)

func writeMetaData(cluster *container.Cluster, project string) {
	// Set up metadata client for saving metadata
	c, err := util.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Writing metadata to: %q", c.Path)
	// Get minNodes and maxNodes counts from default-pool, this is
	// usually the case in tests in Prow
	var minNodes, maxNodes string
	for _, np := range cluster.NodePools {
		if np.Name == "default-pool" {
			if np.Autoscaling != nil {
				minNodes = strconv.FormatInt(np.Autoscaling.MinNodeCount, 10)
				maxNodes = strconv.FormatInt(np.Autoscaling.MaxNodeCount, 10)
			} else {
				log.Printf("DEBUG: nodepool is default-pool but autoscaling is not on: '%+v'", np)
			}
			break
		}
	}

	var e2eRegion, e2eZone string
	geoKey := "E2E:REGION"
	e2eRegion = cluster.Location
	locParts := strings.Split(e2eRegion, "_")
	if len(locParts) == 2 {
		e2eRegion = locParts[0]
		e2eZone = locParts[1]
		geoKey = "E2E:ZONE"
	} else if len(locParts) > 3 {
		log.Fatalf("location %q shouldn't contain more than 1 '_'", cluster.Location)
	}

	for key, val := range map[string]string{
		geoKey:         cluster.Location,
		"E2E:Region":   e2eRegion,
		"E2E:Zone":     e2eZone,
		"E2E:Machine":  cluster.Name,
		"E2E:Version":  cluster.InitialClusterVersion,
		"E2E:MinNodes": minNodes,
		"E2E:MaxNodes": maxNodes,
		"E2E:Project":  project,
	} {
		err = c.Set(key, val)
		if err != nil {
			log.Fatalf("failed saving metadata %q:%q. err: '%v'", key, val, err)
		}
	}
	log.Println("Done writing metadata")
}

func Create(o *options.RequestWrapper) {
	o.Prep()

	gkeClient := clustermanager.GKEClient{}
	clusterOps := gkeClient.Setup(o.Request)
	gkeOps := clusterOps.(*clustermanager.GKECluster)
	if err := gkeOps.Acquire(); err != nil || gkeOps.Cluster == nil {
		log.Fatalf("failed acquiring GKE cluster: '%v'", err)
	}

	// At this point we should have a cluster ready to run test. Need to save
	// metadata so that following flow can understand the context of cluster, as
	// well as for Prow usage later
	// TODO(chaodaiG): this logic may need to be part of clustermanager lib as well
	writeMetaData(gkeOps.Cluster, gkeOps.Project)

	// set up kube config points to cluster
	// TODO(chaodaiG): this probably should also be part of clustermanager lib
	if out, err := common.StandardExec("gcloud", "beta", "container", "clusters", "get-credentials",
		gkeOps.Cluster.Name, "--region", gkeOps.Cluster.Location, "--project", gkeOps.Project); err != nil {
		log.Fatalf("failed connect to cluster: '%v', '%v'", string(out), err)
	}
	if out, err := common.StandardExec("gcloud", "config", "set", "project", gkeOps.Project); err != nil {
		log.Fatalf("failed set gcloud: '%v', '%v'", string(out), err)
	}
}
