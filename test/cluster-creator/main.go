// This is a wrapper of cluster management lib, which can be consumed by command
// line for creating a working cluster
package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	container "google.golang.org/api/container/v1beta1"
	"knative.dev/pkg/testutils/clustermanager"
	"knative.dev/pkg/testutils/common"
	"knative.dev/test-infra/test/metahelper/util"
)

type requestWrapper struct {
	request          clustermanager.GKERequest
	backupRegionsStr string
	addonsStr        string
}

func newRequestWrapper() requestWrapper {
	return requestWrapper{
		request: clustermanager.GKERequest{},
	}
}

func (rw *requestWrapper) prep() {
	if rw.backupRegionsStr != "" {
		rw.request.BackupRegions = strings.Split(rw.backupRegionsStr, ",")
	}
	if rw.addonsStr != "" {
		rw.request.Addons = strings.Split(rw.addonsStr, ",")
	}
}

func writeMetaData(cluster *container.Cluster, project string) {
	// Set up metadata client for saving metadata
	c, err := util.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
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
	log.Println("done writing data to meta.json")
}

func main() {
	rw := newRequestWrapper()
	flag.Int64Var(&rw.request.MinNodes, "min-nodes", 0, "minimal number of nodes")
	flag.Int64Var(&rw.request.MaxNodes, "max-nodes", 0, "maximal number of nodes")
	flag.StringVar(&rw.request.NodeType, "node-type", "", "node type")
	flag.StringVar(&rw.request.Region, "region", "", "GCP region")
	flag.StringVar(&rw.request.Zone, "zone", "", "GCP zone")
	flag.StringVar(&rw.request.Project, "project", "", "GCP project")
	flag.StringVar(&rw.request.ClusterName, "name", "", "cluster name")
	flag.StringVar(&rw.backupRegionsStr, "backup-regions", "", "GCP regions as backup, separated by comma")
	flag.StringVar(&rw.addonsStr, "addons", "", "addons to be added, separated by comma")
	flag.BoolVar(&rw.request.SkipCreation, "skip-creation", false, "should skip creation or not")

	flag.Parse()
	rw.prep()

	gkeClient := clustermanager.GKEClient{}
	clusterOps := gkeClient.Setup(rw.request)
	gkeOps := clusterOps.(*clustermanager.GKECluster)
	if err := gkeOps.Acquire(); err != nil || gkeOps.Cluster == nil {
		log.Fatalf("failed acquiring GKE cluster: '%v'", err)
	}

	// At this point we should have a cluster ready to run test. Need to save
	// metadata so that following flow can understand the context of cluster, as
	// well as for Prow usage later
	// TODO(chaodaiG): this logic may need to be part of clustermanager lib as well
	log.Println("writing data to meta.json")
	writeMetaData(gkeOps.Cluster, *gkeOps.Project)

	// set up kube config points to cluster
	// TODO(chaodaiG): this probably should also be part of clustermanager lib
	if out, err := common.StandardExec("gcloud", "beta", "container", "clusters", "get-credentials",
		gkeOps.Cluster.Name, "--region", gkeOps.Cluster.Location, "--project", *gkeOps.Project); err != nil {
		log.Fatalf("failed connect to cluster: '%v', '%v'", string(out), err)
	}
}
