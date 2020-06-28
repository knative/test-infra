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

// The cleanup tool deletes old images and test clusters from test
// projects.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"knative.dev/test-infra/tools/cleanup/options"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"knative.dev/test-infra/pkg/cmd"
	"knative.dev/test-infra/pkg/gke"
	"knative.dev/test-infra/pkg/helpers"
)

var (
	// Authentication method for using Google Cloud Registry APIs.
	defaultKeychain = authn.DefaultKeychain

	// Alias of remote.Delete for testing purposes.
	remoteDelete = remote.Delete
)

// ResourceDeleter deletes a specific kind of resource in a GCP project.
type ResourceDeleter interface {
	Projects() []string
	Delete(hoursToKeepResource int, concurrentOperations int, dryRun bool) (int, []string)
	DeleteResources(project string, hoursToKeepResource int, dryRun bool) (int, error)
	ShowStats(count int, errors []string)
}

// BaseResourceDeleter implements the base operations of a ResourceDeleter.
type BaseResourceDeleter struct {
	ResourceDeleter
	projects           []string
	deleteResourceFunc func(string, int, bool) (int, error)
}

// ImageDeleter deletes old images in a given registry.
type ImageDeleter struct {
	BaseResourceDeleter
	registry string
}

// GkeClusterDeleter deletes old GKE cluster in a given project.
type GkeClusterDeleter struct {
	BaseResourceDeleter
	gkeClient gke.SDKOperations
}

// NewBaseResourceDeleter returns a brand new BaseResourceDeleter.
func NewBaseResourceDeleter(projects []string) *BaseResourceDeleter {
	deleter := BaseResourceDeleter{projects: projects}
	deleter.deleteResourceFunc = deleter.DeleteResources
	return &deleter
}

// NewImageDeleter returns a brand new ImageDeleter.
func NewImageDeleter(projects []string, registry string, serviceAccount string) (*ImageDeleter, error) {
	var err error
	deleter := ImageDeleter{*NewBaseResourceDeleter(projects), registry}
	deleter.deleteResourceFunc = deleter.DeleteResources
	if serviceAccount != "" {
		// Activate the service account.
		_, err = cmd.RunCommand("gcloud auth activate-service-account --key-file=" + serviceAccount)
		if err != nil {
			if cmdErr, ok := err.(*cmd.CommandLineError); ok {
				err = fmt.Errorf("cannot activate service account:\n%s", cmdErr.ErrorOutput)
			}
		}
	}
	return &deleter, err
}

// NewGkeClusterDeleter returns a brand new GkeClusterDeleter.
func NewGkeClusterDeleter(projects []string, serviceAccount string) (*GkeClusterDeleter, error) {
	opts := make([]option.ClientOption, 0)
	if serviceAccount != "" {
		// Create GKE client with specific credentials.
		opts = append(opts, option.WithCredentialsFile(serviceAccount))
	}
	gkeClient, err := gke.NewSDKClient(opts...)
	if err != nil {
		err = errors.Wrapf(err, "cannot create GKE SDK client")
	}
	deleter := GkeClusterDeleter{*NewBaseResourceDeleter(projects), gkeClient}
	deleter.deleteResourceFunc = deleter.DeleteResources
	return &deleter, err
}

// selectProjects returns the list of projects to iterate over.
func selectProjects(projects []string, resourceFiles []string, regex string) ([]string, error) {
	// Sanity check flags
	if len(projects) == 0 && len(resourceFiles) == 0 {
		return nil, errors.New("neither project nor resource file provided")
	}

	if len(projects) != 0 && len(resourceFiles) > 0 {
		return nil, errors.New("provided both project and resource file")
	}
	// --project used, just return it.
	if len(projects) != 0 {
		log.Printf("Iterating over projects %v", projects)
		return projects, nil
	}

	// Otherwise, read the resource file and extract the project names.
	return fromResourceFiles(resourceFiles, regex)
}

// selectProjects returns the list of projects to iterate over.
func fromResourceFiles(resourceFiles []string, regex string) ([]string, error) {
	var projects []string
	projectRegex, err := regexp.Compile(regex)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid regular expression %q", regex)
	}
	for _, resourceFile := range resourceFiles {
		content, err := ioutil.ReadFile(resourceFile)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read file %q", resourceFile)
		}
		log.Printf("Iterating over projects defined in %q, matching %q", resourceFile, regex)
		var count int
		for _, line := range strings.Split(string(content), "\n") {
			if len(line) > 0 {
				if p := projectRegex.Find([]byte(line)); p != nil {
					log.Printf("\t- %s", string(p))
					projects = append(projects, string(p))
					count++
				}
			}
		}
		log.Printf("Found %d projects", count)
	}
	if len(projects) == 0 {
		return nil, fmt.Errorf("no project found in '%v' matching %q", resourceFiles, regex)
	}
	return projects, nil
}

// deleteImage deletes a single image pointed by the given reference.
func (d *ImageDeleter) deleteImage(ref string) error {
	image, err := name.ParseReference(ref)
	if err != nil {
		return errors.Wrapf(err, "failed to parse reference %q", ref)
	}
	if err := remoteDelete(image, remote.WithAuthFromKeychain(defaultKeychain)); err != nil {
		return errors.Wrapf(err, "failed to delete %q", image)
	}
	return nil
}

// DeleteResources deletes old clusters from a given project.
func (d *GkeClusterDeleter) DeleteResources(project string, hoursToKeepResource int, dryRun bool) (int, error) {
	before := time.Now().Add(-time.Hour * time.Duration(hoursToKeepResource))
	// TODO(adrcunha): Consider exposing https://github.com/knative/pkg/blob/6d806b998379948bd0107d77bcd831e2bdb4f3cb/testutils/clustermanager/e2e-tests/gke.go#L281
	if project == "knative-tests" {
		return 0, fmt.Errorf("cleaning up %q is forbidden", project)
	}
	// List clusters, delete those created before the given timestamp.
	clusters, err := d.gkeClient.ListClustersInProject(project)
	if err != nil {
		return 0, errors.Wrapf(err, "error listing clusters in %q, maybe try 'gcloud auth application-default login'", project)
	}
	count := 0
	for _, cluster := range clusters {
		creation, err := time.Parse(time.RFC3339, cluster.CreateTime)
		if err != nil {
			return count, errors.Wrapf(err, "error getting creation time for cluster %q", cluster.Name)
		}
		age := int(time.Since(creation).Hours())
		fullClusterName := project + "/" + cluster.Name
		log.Printf("%s is %d hours old", fullClusterName, age)
		if creation.Before(before) {
			if err := helpers.Run(fmt.Sprintf("Deleting %q", fullClusterName), func() error {
				region, zone := gke.RegionZoneFromLoc(cluster.Location)
				if err := d.gkeClient.DeleteCluster(project, region, zone, cluster.Name); err != nil {
					return errors.Wrapf(err, "error deleting cluster %q in project %q", cluster.Name, project)
				}
				count++
				return nil
			}, dryRun); err != nil {
				return count, err
			}
		}
	}
	return count, nil
}

// DeleteResources deletes old docker images from a given project.
func (d *ImageDeleter) DeleteResources(project string, hoursToKeepResource int, dryRun bool) (int, error) {
	before := time.Now().Add(-time.Hour * time.Duration(hoursToKeepResource))
	repoRoot := d.registry + "/" + project
	// TODO(adrcunha): This should be a helper function, like https://github.com/knative/pkg/blob/6d806b998379948bd0107d77bcd831e2bdb4f3cb/testutils/clustermanager/e2e-tests/gke.go#L281
	if repoRoot == "gcr.io/knative-releases" || repoRoot == "gcr.io/knative-nightly" {
		return 0, fmt.Errorf("cleaning up %q is forbidden", repoRoot)
	}
	gcrrepo, err := name.NewRepository(repoRoot)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot open registry %q", repoRoot)
	}
	count := 0
	// Walk down the registry, checking all images and deleting the old ones.
	return count, google.Walk(gcrrepo, func(repo name.Repository, tags *google.Tags, err error) error {
		// If we got an error, just return it, there's nothing to do here.
		if tags == nil || err != nil {
			if err == nil {
				return fmt.Errorf("unexpected nil tags for %q", repo)
			}
			return errors.Wrapf(err, "cannot walk down GCR %q", repo.String())
		}
		for k, m := range tags.Manifests {
			ref := repo.String() + "@" + k
			age := int(time.Since(m.Uploaded).Hours() / 24)
			log.Printf("%q is %d days old (uploaded on %s)", ref, age, m.Uploaded)
			if m.Uploaded.Before(before) {
				if err := helpers.Run(fmt.Sprintf("Deleting %q", ref), func() error {
					// Delete all tags first, otherwise the image can't be deleted.
					for _, tag := range m.Tags {
						if err := d.deleteImage(repo.String() + ":" + tag); err != nil {
							return err
						}
					}
					if err := d.deleteImage(ref); err != nil {
						return err
					}
					count++
					return nil
				}, dryRun); err != nil {
					return err
				}
			}
		}
		return nil
	}, google.WithAuthFromKeychain(defaultKeychain))
}

// Projects returns the projects that should be cleaned up by a ResourceDeleter.
func (d *BaseResourceDeleter) Projects() []string {
	return d.projects
}

// DeleteResources base method that does nothing, as it must be overridden.
func (d *BaseResourceDeleter) DeleteResources(project string, hoursToKeepResource int, dryRun bool) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

// Delete call DeleteResource in parallel, one for each given project.
func (d *BaseResourceDeleter) Delete(hoursToKeepResource int, concurrentOperations int, dryRun bool) (int, []string) {
	// Locks for the concurrent tasks.
	var wg sync.WaitGroup

	// Blocking channel to keep concurrency under control.
	concurrencyChan := make(chan struct{}, concurrentOperations)

	projects := d.Projects()

	// Channel to hold errors.
	errorChan := make(chan error, len(projects))

	var count int32
	for i := range projects {
		wg.Add(1)
		go func(project string) {
			defer wg.Done()
			// Block processing if concurrency reached the limit.
			concurrencyChan <- struct{}{}
			// Do not process if previous invocations failed. This prevents a large
			// build-up of failed requests and rate limit exceeding (e.g. bad auth).
			if len(errorChan) == 0 {
				c, err := d.deleteResourceFunc(project, hoursToKeepResource, dryRun)
				// Update counter and errors list.
				atomic.AddInt32(&count, int32(c))
				if err != nil {
					errorChan <- err
				}
			}
			<-concurrencyChan
		}(projects[i])
	}
	wg.Wait()
	close(errorChan)
	close(concurrencyChan)
	// Extract the error strings from the map of errors.
	// For testing purposes, sort them to keep order constant.
	errStrings := make([]string, 0)
	for e := range errorChan {
		unique := true
		for _, s := range errStrings {
			if s == e.Error() {
				unique = false
				break
			}
		}
		if unique {
			errStrings = append(errStrings, e.Error())
		}
	}
	sort.Strings(errStrings)
	return int(count), errStrings
}

// ShowStats simply shows the number of resources deleted, and any errors.
func (d *BaseResourceDeleter) ShowStats(count int, errors []string) {
	log.Printf("%d resources deleted", count)
	if len(errors) > 0 {
		log.Printf("%d errors occurred: %s", len(errors), strings.Join(errors, ", "))
	}
}

// cleanup parses flags, run the operations and returns the status.
func cleanup(o options.Options) error {
	if o.DryRun {
		log.Println("-- Running in dry-run mode, no resource deletion --")
	}

	if !strings.HasSuffix(o.Registry, "gcr.io") {
		return fmt.Errorf("currently only GCR is supported")
	}

	var projects []string
	var err error
	if projects, err = selectProjects(o.Project, o.ProjectResourceYaml, o.ReProjectName); err != nil {
		return err
	}

	start := time.Now()

	var deleter ResourceDeleter
	if o.DaysToKeepImages >= 0 {
		if deleter, err = NewImageDeleter(projects, o.Registry, o.ServiceAccount); err != nil {
			return err
		}
		log.Println("Removing images that are:")
		log.Printf("- older than %d days", o.DaysToKeepImages)
		deleter.ShowStats(deleter.Delete(o.DaysToKeepImages*24, o.ConcurrentOperations, o.DryRun))
	}

	if o.HoursToKeepClusters >= 0 {
		if deleter, err = NewGkeClusterDeleter(projects, o.ServiceAccount); err != nil {
			return err
		}
		log.Println("Removing clusters that are:")
		log.Printf("- older than %d hours", o.HoursToKeepClusters)
		deleter.ShowStats(deleter.Delete(o.HoursToKeepClusters, o.ConcurrentOperations, o.DryRun))
	}

	log.Printf("All operations finished in %s", time.Now().Sub(start))
	return nil
}

// main is the script entry point.
func main() {
	var o options.Options
	o.AddOptions()
	flag.Parse()

	if err := cleanup(o); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
