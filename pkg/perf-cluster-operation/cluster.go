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

package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"knative.dev/pkg/test/helpers"
	"knative.dev/pkg/testutils/gke"

	container "google.golang.org/api/container/v1beta1"
)

const (
	// the maximum retry times if there is an error in cluster operation
	retryTimes = 3

	// known cluster status
	statusProvisioning = "PROVISIONING"
	statusRunning      = "RUNNING"
	statusStopping     = "STOPPING"
)

type gkeClient struct {
	ops gke.SDKOperations
}

func newClient() (*gkeClient, error) {
	operations, err := gke.NewSDKClient()
	if err != nil {
		return nil, fmt.Errorf("failed to set up GKE client: %v", err)
	}

	client := &gkeClient{
		ops: operations,
	}
	return client, nil
}

// recreateClusters will delete and recreate the existing clusters
func (gc *gkeClient) recreateClusters(gcpProject, repo, benchmarkRoot string) error {
	handleExistingCluster := func(cluster container.Cluster, clusterConfigs map[string]ClusterConfig) error {
		// always delete the cluster, even if the cluster config is unchanged
		return gc.handleExistingClusterHelper(gcpProject, cluster, clusterConfigs, false)
	}
	handleNewClusterConfig := func(clusterName string, clusterConfig ClusterConfig) error {
		// for now, do nothing to the new cluster config
		return nil
	}
	return gc.processClusters(gcpProject, repo, benchmarkRoot, handleExistingCluster, handleNewClusterConfig)
}

// reconcileClusters will reconcile all clusters to make them consistent with the benchmarks' cluster configs
//
// There can be 4 scenarios:
// 1. If the benchmark's cluster config is unchanged, do nothing
// 2. If the benchmark's config is changed, delete the old cluster and create a new one with the new config
// 3. If the benchmark is renamed, delete the old cluster and create a new one with the new name
// 4. If the benchmark is deleted, delete the corresponding cluster
func (gc *gkeClient) reconcileClusters(gcpProject, repo, benchmarkRoot string) error {
	handleExistingCluster := func(cluster container.Cluster, clusterConfigs map[string]ClusterConfig) error {
		// retain the cluster, if the cluster config is unchanged
		return gc.handleExistingClusterHelper(gcpProject, cluster, clusterConfigs, true)
	}
	handleNewClusterConfig := func(clusterName string, clusterConfig ClusterConfig) error {
		// create a new cluster with the new cluster config
		return gc.createClusterWithRetries(gcpProject, clusterName, clusterConfig)
	}
	return gc.processClusters(gcpProject, repo, benchmarkRoot, handleExistingCluster, handleNewClusterConfig)
}

// processClusters will process existing clusters and configs for new clusters,
// with the corresponding functions provided by callers.
func (gc *gkeClient) processClusters(
	gcpProject, repo, benchmarkRoot string,
	handleExistingCluster func(cluster container.Cluster, clusterConfigs map[string]ClusterConfig) error,
	handleNewClusterConfig func(name string, config ClusterConfig) error,
) error {
	curtClusters, err := gc.listClustersForRepo(gcpProject, repo)
	if err != nil {
		return fmt.Errorf("failed getting clusters for the repo %q: %v", repo, err)
	}
	clusterConfigs, err := benchmarkClusters(repo, benchmarkRoot)
	if err != nil {
		return fmt.Errorf("failed getting cluster configs for benchmarks in repo %q: %v", repo, err)
	}

	errCh := make(chan error)
	wg := sync.WaitGroup{}
	// handle all existing clusters
	for i := range curtClusters {
		wg.Add(1)
		cluster := curtClusters[i]
		go func() {
			defer wg.Done()
			if err := handleExistingCluster(cluster, clusterConfigs); err != nil {
				errCh <- fmt.Errorf("failed handling cluster %v: %v", cluster, err)
			}
		}()
		// remove the cluster from clusterConfigs as it's already been handled
		delete(clusterConfigs, cluster.Name)
	}

	// handle all other cluster configs
	for name, config := range clusterConfigs {
		wg.Add(1)
		newName, newConfig := name, config
		go func() {
			defer wg.Done()
			if err := handleNewClusterConfig(newName, newConfig); err != nil {
				errCh <- fmt.Errorf("failed handling new cluster config %v: %v", newConfig, err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	errs := make([]error, 0)
	for err := range errCh {
		errs = append(errs, err)
	}

	return helpers.CombineErrors(errs)
}

// handleExistingClusterHelper is a helper function for handling an existing cluster.
func (gc *gkeClient) handleExistingClusterHelper(
	gcpProject string, cluster container.Cluster, clusterConfigs map[string]ClusterConfig,
	retainIfUnchanged bool,
) error {
	// if the cluster is currently being created or deleted, return directly as other jobs will handle it properly
	if cluster.Status == statusProvisioning || cluster.Status == statusStopping {
		log.Printf("cluster %q is being handled by other jobs, skip it", cluster.Name)
		return nil
	}
	// if retainIfUnchanged is set to true, and the cluster config does not change, do nothing
	// TODO(chizhg): also check the addons config
	config, configExists := clusterConfigs[cluster.Name]
	if retainIfUnchanged && configExists &&
		cluster.CurrentNodeCount == config.NodeCount && cluster.Location == config.Location {
		log.Printf("cluster config is unchanged for %q, skip it", cluster.Name)
		return nil
	}

	if err := gc.deleteClusterWithRetries(gcpProject, cluster); err != nil {
		return fmt.Errorf("failed deleting cluster %q in %q: %v", cluster.Name, cluster.Location, err)
	}
	if configExists {
		return gc.createClusterWithRetries(gcpProject, cluster.Name, config)
	}
	return nil
}

// listClustersForRepo will list all the clusters under the gcpProject that belong to the given repo.
func (gc *gkeClient) listClustersForRepo(gcpProject, repo string) ([]container.Cluster, error) {
	allClusters, err := gc.ops.ListClustersInProject(gcpProject)
	if err != nil {
		return nil, fmt.Errorf("failed listing clusters in project %q: %v", gcpProject, err)
	}

	clusters := make([]container.Cluster, 0)
	for _, cluster := range allClusters {
		if clusterBelongsToRepo(cluster.Name, repo) {
			clusters = append(clusters, *cluster)
		}
	}
	return clusters, nil
}

// deleteClusterWithRetries will delete the given cluster,
// and retry for a maximum of retryTimes if there is an error.
// TODO(chizhg): maybe move it to clustermanager library.
func (gc *gkeClient) deleteClusterWithRetries(gcpProject string, cluster container.Cluster) error {
	region, zone := gke.RegionZoneFromLoc(cluster.Location)
	var err error
	for i := 0; i < retryTimes; i++ {
		if err = gc.ops.DeleteCluster(gcpProject, region, zone, cluster.Name); err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf(
			"failed deleting cluster %q in %q after retrying %d times: %v",
			cluster.Name, gke.GetClusterLocation(region, zone), retryTimes, err)
	}

	return nil
}

// createClusterWithRetries will create a new cluster with the given config,
// and retry for a maximum of retryTimes if there is an error.
// TODO(chizhg): maybe move it to clustermanager library.
func (gc *gkeClient) createClusterWithRetries(gcpProject, name string, config ClusterConfig) error {
	req := &gke.Request{
		ClusterName: name,
		MinNodes:    config.NodeCount,
		MaxNodes:    config.NodeCount,
		NodeType:    config.NodeType,
		Addons:      strings.Split(config.Addons, ","),
	}
	creq, err := gke.NewCreateClusterRequest(req)
	if err != nil {
		return fmt.Errorf("cannot create cluster with request %v: %v", req, err)
	}

	region, zone := gke.RegionZoneFromLoc(config.Location)
	for i := 0; i < retryTimes; i++ {
		// TODO(chizhg): retry with different requests, based on the error type
		if err = gc.ops.CreateCluster(gcpProject, region, zone, creq); err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf(
			"failed creating cluster %q in %q after retrying %d times: %v",
			name, config.Location, retryTimes, err)
	}

	return nil
}
