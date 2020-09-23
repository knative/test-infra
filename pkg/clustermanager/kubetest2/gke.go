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

package kubetest2

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"knative.dev/test-infra/pkg/cmd"
	"knative.dev/test-infra/pkg/metautil"
	"knative.dev/test-infra/pkg/prow"
)

const (
	createCommandTmpl = "%s container clusters create --quiet --enable-autoscaling --min-nodes=%d --max-nodes=%d " +
		"--scopes=%s"
	boskosAcquireDefaultTimeoutSeconds = 1200
)

var (
	baseKubetest2Flags = []string{"gke", "--ignore-gcp-ssh-key=true", "--up"}

	// If one of the error patterns below is matched, it would be recommended to
	// retry creating the cluster in a different region.
	// - stockout (https://github.com/knative/test-infra/issues/592)
	retryableCreationErrors = []*regexp.Regexp{
		regexp.MustCompile(`.*does not have enough resources available to fulfill.*`),
		regexp.MustCompile(`.*only \d+ nodes out of \d+ have registered; this is likely due to Nodes failing to start correctly.*`),
	}
)

// GKEClusterConfig are the supported configurations for creating a GKE cluster.
type GKEClusterConfig struct {
	GCPServiceAccount string
	GCPProjectID      string

	BoskosAcquireTimeoutSeconds int

	Environment  string
	CommandGroup string

	Name                              string
	Region                            string
	BackupRegions                     []string
	Machine                           string
	MinNodes                          int
	MaxNodes                          int
	Network                           string
	ReleaseChannel                    string
	Version                           string
	Scopes                            string
	Addons                            string
	EnableWorkloadIdentity            bool
	PrivateClusterAccessLevel         string
	PrivateClusterMasterIPSubnetRange string
	PrivateClusterMasterIPSubnetMask  string

	ExtraGcloudFlags string
}

// Run will run the `kubetest2 gke` command with the provided parameters,
// it will also handle the logic that is only used for Knative integration testing, like retrying cluster creation.
func Run(opts *Options, cc *GKEClusterConfig) error {
	createCommand := fmt.Sprintf(createCommandTmpl, cc.CommandGroup, cc.MinNodes, cc.MaxNodes, cc.Scopes)
	// --cluster-version must be a valid version number when --release-channel is not empty,
	// so normally we do not use them together.
	if cc.ReleaseChannel != "" {
		createCommand += " --release-channel=" + cc.ReleaseChannel
	} else {
		createCommand += " --cluster-version=" + cc.Version
	}
	if cc.Addons != "" {
		createCommand += " --addons=" + cc.Addons
	}
	if cc.ExtraGcloudFlags != "" {
		createCommand += " " + cc.ExtraGcloudFlags
	}
	kubetest2Flags := append(baseKubetest2Flags, "--create-command="+createCommand)

	kubetest2Flags = append(kubetest2Flags, "--cluster-name="+cc.Name, "--environment="+cc.Environment,
		"--num-nodes="+strconv.Itoa(cc.MinNodes), "--machine-type="+cc.Machine, "--network="+cc.Network)
	if cc.GCPServiceAccount != "" {
		kubetest2Flags = append(kubetest2Flags, "--gcp-service-account="+cc.GCPServiceAccount)
	}
	if cc.EnableWorkloadIdentity {
		kubetest2Flags = append(kubetest2Flags, "--enable-workload-identity")
	}

	if prow.IsCI() && cc.GCPProjectID == "" {
		log.Println("Will use boskos to provision the GCP project")
		timeout := cc.BoskosAcquireTimeoutSeconds
		if timeout == 0 {
			timeout = boskosAcquireDefaultTimeoutSeconds
		}
		kubetest2Flags = append(kubetest2Flags, "--boskos-acquire-timeout-seconds="+strconv.Itoa(timeout))
	} else {
		if cc.GCPProjectID == "" {
			return errors.New("GCP project must be provided in non-CI environment")
		}
		log.Printf("Will use the GCP project %q for creating the cluster", cc.GCPProjectID)
		kubetest2Flags = append(kubetest2Flags, "--project="+cc.GCPProjectID)
	}

	return createGKEClusterWithRetries(kubetest2Flags, opts, cc)
}

func createGKEClusterWithRetries(kubetest2Flags []string, opts *Options, cc *GKEClusterConfig) error {
	var err error
	regions := append([]string{cc.Region}, cc.BackupRegions...)
	for i, region := range regions {
		flags := append(kubetest2Flags, "--region="+region)
		if cc.PrivateClusterAccessLevel != "" {
			flags = append(flags, "--private-cluster-access-level="+cc.PrivateClusterAccessLevel)
			masterIPRange := fmt.Sprintf("%s.%d/%s", cc.PrivateClusterMasterIPSubnetRange, i, cc.PrivateClusterMasterIPSubnetMask)
			flags = append(flags, "--private-cluster-master-ip-range="+masterIPRange)
		}

		if opts.TestCommand != "" {
			flags = append(flags, "--test=exec", "--")
			flags = append(flags, strings.Split(opts.TestCommand, " ")...)
		}

		log.Printf("Running kubetest2 with flags: %q", flags)
		command := exec.Command("kubetest2", flags...)
		var out string
		out, err = runWithOutput(command)
		if err != nil {
			if isRetryableCreationError(out) {
				log.Print("Cluster creation fails due to unpredictable reasons, will retry creating with different args")
				continue
			} else {
				return err
			}
		} else {
			// Only save the metadata if it's in CI environment and meta data is asked to be saved.
			if prow.IsCI() && opts.SaveMetaData {
				saveMetaData(cc, region)
			}
			break
		}
	}

	return err
}

// isRetryableCreationError determines if cluster creation should be retried based on the error message.
func isRetryableCreationError(errMsg string) bool {
	for _, regx := range retryableCreationErrors {
		if regx.MatchString(errMsg) {
			return true
		}
	}
	return false
}

// saveMetaData will save the metadata with best effort.
func saveMetaData(cc *GKEClusterConfig, region string) {
	cli, err := metautil.NewClient("")
	if err != nil {
		log.Printf("error creating the metautil client: %v", err)
		return
	}
	cv, err := cmd.RunCommand("kubectl version --short=true")
	if err != nil {
		log.Printf("error getting the cluster version: %v", err)
		return
	}

	// Set the metadata with best effort.
	cli.Set("E2E:Provider", "gke")
	cli.Set("E2E:Region", region)
	cli.Set("E2E:Machine", cc.Machine)
	cli.Set("E2E:Version", cv)
	cli.Set("E2E:MinNodes", strconv.Itoa(cc.MinNodes))
	cli.Set("E2E:MaxNodes", strconv.Itoa(cc.MaxNodes))
}
