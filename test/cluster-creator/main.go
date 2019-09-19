// This is just an exmaple of creating most basic cluster for integration test
// purpose
package main

import (
	"log"
	"strconv"

	"knative.dev/eventing-operator/test/metahelper/util"

	"knative.dev/pkg/testutils/clustermanager"
)

func main() {
	gkeClient := clustermanager.GKEClient{}
	clusterOps := gkeClient.Setup(nil, nil, nil, nil, nil, nil, []string{})
	gkeOps := clusterOps.(*clustermanager.GKECluster)
	if err := gkeOps.Initialize(); err != nil {
		log.Fatalf("failed initializing GKE Client: '%v'", err)
	}
	if err := gkeOps.Acquire(); err != nil {
		log.Fatalf("failed acquire cluster: '%v'", err)
	}
	// At this point we should have a cluster ready to run test. Need to save
	// metadata so that following flow can understand the context of cluster, as
	// well as for Prow usage later
	// TODO(chaodaiG): this logic may need to be part of clustermanager lib
	c, err := util.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	// Get minNodes and maxNodes counts from 1st available node pool, this is
	// usually the case in tests in Prow
	var minNodes, maxNodes string
	if len(gkeOps.Cluster.NodePools) > 0 {
		minNodes = strconv.FormatInt(gkeOps.Cluster.NodePools[0].Autoscaling.MinNodeCount, 10)
		maxNodes = strconv.FormatInt(gkeOps.Cluster.NodePools[0].Autoscaling.MaxNodeCount, 10)
	}
	geoKey := "E2E:REGION"
	if gkeOps.Request.Zone != "" {
		geoKey = "E2E:ZONE"
	}
	for key, val := range map[string]string{
		geoKey:         gkeOps.Cluster.Location,
		"E2E:Machine":  gkeOps.Cluster.Name,
		"E2E:Version":  gkeOps.Cluster.InitialClusterVersion,
		"E2E:MinNodes": minNodes,
		"E2E:MaxNodes": maxNodes,
	} {
		err = c.Set(key, val)
		if err != nil {
			log.Fatalf("failed saving metadata %q:%q. err: '%v'", key, val, err)
		}
	}
}
