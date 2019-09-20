// This is a wrapper of cluster management lib, which can be consumed by command
// line for deleting a working cluster and/or releasing Boskos project used for
// testing
package main

import (
	"log"

	"knative.dev/pkg/testutils/clustermanager"
)

func main() {
	r := clustermanager.GKERequest{
		NeedsCleanup: true,
		SkipCreation: true,
	}
	gkeClient := clustermanager.GKEClient{}
	clusterOps := gkeClient.Setup(r)
	gkeOps := clusterOps.(*clustermanager.GKECluster)
	if err := gkeOps.Acquire(); err != nil || gkeOps.Cluster == nil {
		log.Fatal("Failed identifying cluster for cleanup")
	}
	log.Printf("Identified project %q and cluster %q for removal", *gkeOps.Project, gkeOps.Cluster.Name)
	// Don't wait for delete
	if err := gkeOps.Delete(); err != nil {
		log.Fatalf("Failed deleting cluster: '%v'", err)
	}
}
