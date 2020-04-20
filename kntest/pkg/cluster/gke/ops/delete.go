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

package ops

import (
	"fmt"
	"log"

	"knative.dev/pkg/test/cmd"
	clm "knative.dev/pkg/testutils/clustermanager/e2e-tests"
)

// Delete deletes a GKE cluster
func (o *RequestWrapper) Delete() error {
	o.Request.SkipCreation = true

	gkeClient := clm.GKEClient{}
	clusterOps := gkeClient.Setup(o.Request)
	gkeOps := clusterOps.(*clm.GKECluster)
	if err := gkeOps.Acquire(); err != nil || gkeOps.Cluster == nil {
		return fmt.Errorf("failed identifying cluster for cleanup: '%w'", err)
	}
	log.Printf("Identified project %q and cluster %q for removal", gkeOps.Project, gkeOps.Cluster.Name)
	var err error
	if err = gkeOps.Delete(); err != nil {
		return fmt.Errorf("failed deleting cluster: '%w'", err)
	}
	// Unset context with best effort.
	// The commands will try to unset the current context and delete it from kubeconfig.
	cc, _ := cmd.RunCommand("kubectl config current-context")
	if _, err := cmd.RunCommand("kubectl config unset current-context"); err != nil {
		cmd.RunCommand("kubectl config unset contexts." + cc)
	}

	return nil
}
