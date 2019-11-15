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

// The cleanup tool deletes old GCR images and test clusters from test
// projects.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"knative.dev/test-infra/shared/common"

	"knative.dev/pkg/test/gke"
)

// resourceDeleter is a function that deletes a particular type of resource, and return how many resources were deleted.
type resourceDeleter func(string, time.Time) (int, error)

var (
	// Command-line flags.
	projectResourceYaml  = flag.String("project-resource-yaml", "", "Resources file containing the names of the projects to be cleaned up.")
	reProjectName        = flag.String("re-project-name", "knative-boskos-[a-zA-Z0-9]+", "Regular expression for filtering project names from the resources file.")
	daysToKeepImages     = flag.Int("days-to-keep-images", 365, "Images older than this amount of days will be deleted (defaults to 1 year, -1 means 'forever').")
	hoursToKeepClusters  = flag.Int("hours-to-keep-clusters", 720, "Clusters older than this amount of hours will be deleted (defautls to 1 month, -1 means 'forever').")
	project              = flag.String("project", "", "Project to be cleaned up.")
	gcr                  = flag.String("gcr", "gcr.io", "The GCR hostname to use.")
	serviceAccount       = flag.String("service-account", "", "Specify the key file of the service account to use.")
	dryRun               = flag.Bool("dry-run", false, "Performs a dry run for all deletion functions.")
	concurrentOperations = flag.Int("concurrent-operations", 10, "How many deletion operations to run concurrently.")

	// Authentication method for using Google Cloud Registry APIs.
	auther = authn.DefaultKeychain

	// Client for deleting GKE clusters.
	gkeClient gke.SDKOperations
)

// selectProjects returns the list of projects to iterate over.
func selectProjects(project, resourceFile, regex string) ([]string, error) {
	// --project used, just return it.
	if project != "" {
		log.Printf("Iterating over projects [%s]", project)
		return []string{project}, nil
	}
	// Otherwise, read the resource file and extract the project names.
	projectRegex, err := regexp.Compile(regex)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid regular expression %q", regex)
	}
	content, err := ioutil.ReadFile(resourceFile)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read file %q", resourceFile)
	}
	projects := make([]string, 0)
	for _, line := range strings.Split(string(content), "\n") {
		if len(line) > 0 {
			if p := projectRegex.Find([]byte(line)); p != nil {
				projects = append(projects, string(p))
			}
		}
	}
	if len(projects) == 0 {
		return nil, fmt.Errorf("no project found in %q matching %q", resourceFile, regex)
	}
	log.Printf("Iterating over projects defined in %q, matching %q", resourceFile, regex)
	return projects, nil
}

// deleteImage deletes a single GCR image pointed by the given reference.
func deleteImage(ref string) error {
	image, err := name.ParseReference(ref)
	if err != nil {
		return errors.Wrapf(err, "failed to parse reference %q", ref)
	}
	if err := remote.Delete(image, remote.WithAuthFromKeychain(auther)); err != nil {
		return errors.Wrapf(err, "failed to delete %s", image)
	}
	return nil
}

// deleteClusters deletes old clusters from a given project.
func deleteClusters(project string, before time.Time) (int, error) {
	// TODO(adrcunha): Consider exposing https://github.com/knative/pkg/blob/6d806b998379948bd0107d77bcd831e2bdb4f3cb/testutils/clustermanager/e2e-tests/gke.go#L281
	if project == "knative-tests" {
		return 0, fmt.Errorf("cleaning up %q is forbidden", project)
	}
	// List clusters, delete those created before the given timestamp.
	clusters, err := gkeClient.ListClustersInProject(project)
	if err != nil {
		return 0, errors.Wrap(err, "error listing clusters in %q, maybe try 'gcloud auth application-default login'")
	}
	count := 0
	for _, cluster := range clusters {
		creation, err := time.Parse(time.RFC3339, cluster.CreateTime)
		if err != nil {
			return count, errors.Wrapf(err, "error getting creation time for cluster %q", cluster.Name)
		}
		age := int(time.Since(creation).Hours())
		log.Printf("%s/%s is %d hours old", project, cluster.Name, age)
		if creation.Before(before) {
			if *dryRun {
				log.Printf("[DRY RUN] delete %q", project+"/"+cluster.Name)
			} else {
				region, zone := gke.RegionZoneFromLoc(cluster.Location)
				if err := gkeClient.DeleteCluster(project, region, zone, cluster.Name); err != nil {
					return count, errors.Wrapf(err, "error deleting cluster %q in project %q", cluster.Name, project)
				}
				count++
			}
		}
	}
	return count, nil
}

// deleteImages deletes old docker images from a given project.
func deleteImages(project string, before time.Time) (int, error) {
	repoRoot := *gcr + "/" + project
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
		if tags == nil {
			return err
		}
		for k, m := range tags.Manifests {
			ref := repo.String() + "@" + k
			age := int(time.Since(m.Uploaded).Hours() / 24)
			log.Printf("%q is %d days old", ref, age)
			if m.Uploaded.Before(before) {
				if *dryRun {
					log.Printf("[DRY RUN] delete %q", ref)
				} else {
					// Delete all tags first, otherwise the image can't be deleted.
					for _, tag := range m.Tags {
						if err := deleteImage(repo.String() + ":" + tag); err != nil {
							return err
						}
					}
					if err := deleteImage(ref); err != nil {
						return err
					}
					log.Printf("Deleted %q (uploaded on %s)", ref, m.Uploaded)
					count++
				}
			}
		}
		return nil
	}, google.WithAuthFromKeychain(auther))
}

// deleteResources call resourceDeleter functions in parallel, one for each given project.
func deleteResources(projects []string, before time.Time, deleteFunc resourceDeleter) (int, []string) {
	// Map of errors returned from functions, to avoid repeated entries.
	var errs = make(map[string]error)

	// Locks for the concurrent tasks.
	var errsLock sync.RWMutex
	var countLock sync.RWMutex
	var wg sync.WaitGroup

	// Blocking channel to keep concurrency under control.
	concurrencyChan := make(chan struct{}, *concurrentOperations)
	defer close(concurrencyChan)

	count := 0
	for i := range projects {
		wg.Add(1)
		go func(project string) {
			defer wg.Done()
			// Block processing if concurrency reached the limit.
			concurrencyChan <- struct{}{}
			// Do not process if previous invocations failed. This prevents a large
			// build-up of failed requests and rate limit exceeding (e.g. bad auth).
			errsLock.RLock()
			failed := len(errs) > 0
			errsLock.RUnlock()
			if failed {
				return
			}
			c, err := deleteFunc(project, before)
			// Update counter and errors list.
			countLock.Lock()
			count += c
			countLock.Unlock()
			if err != nil {
				cause := errors.Cause(err).Error()
				errsLock.Lock()
				if _, ok := errs[cause]; !ok {
					errs[cause] = err
				}
				errsLock.Unlock()
			}
			<-concurrencyChan
		}(projects[i])
	}
	wg.Wait()
	// Extract the error strings from the map of errors.
	errStrings := make([]string, 0)
	for i := range errs {
		errStrings = append(errStrings, errs[i].Error())
	}
	return count, errStrings
}

// showStats displays basic stats based on the return values from deleteResources.
func showStats(count int, errs []string) {
	log.Printf("%d resources deleted", count)
	if len(errs) > 0 {
		log.Printf("%d errors occurred: %s", len(errs), strings.Join(errs, ", "))
	}
}

// cleanup performs the main operations and return any errors.
func cleanup() error {
	// Sanity check flags
	if *project == "" && *projectResourceYaml == "" {
		return errors.New("neither project nor resource file provided")
	}

	if *project != "" && *projectResourceYaml != "" {
		return errors.New("provided both project and resource file")
	}

	if *dryRun {
		log.Println("-- Running in dry-run mode, no resource deletion --")
	}

	// Perform authentication if necessary, and create required clients.
	opts := make([]option.ClientOption, 0)
	if *serviceAccount != "" {
		// Activate the service account.
		if _, output, err := common.ExecCommand("gcloud", "auth", "activate-service-account", "--key-file="+*serviceAccount); err != nil {
			return fmt.Errorf("cannot activate service account:\n%s", output)
		}
		// Create GKE client with specific credentials.
		opts = append(opts, option.WithCredentialsFile(*serviceAccount))
	}
	var err error
	if gkeClient, err = gke.NewSDKClient(opts...); err != nil {
		return errors.Wrapf(err, "cannot create GKE SDK client")
	}

	// Select projects, perform operations.
	var projects []string
	if projects, err = selectProjects(*project, *projectResourceYaml, *reProjectName); err != nil {
		return err
	}
	start := time.Now()

	if *daysToKeepImages >= 0 {
		log.Println("Removing images that are:")
		log.Printf("- older than %d days", *daysToKeepImages)
		showStats(deleteResources(projects, start.Add(-24*time.Hour*time.Duration(*daysToKeepImages)), deleteImages))
	}

	if *hoursToKeepClusters >= 0 {
		log.Println("Removing clusters that are:")
		log.Printf("- older than %d hours", *hoursToKeepClusters)
		showStats(deleteResources(projects, start.Add(-time.Hour*time.Duration(*hoursToKeepClusters)), deleteClusters))
	}

	log.Printf("All operations finished in %s", time.Now().Sub(start))

	return nil
}

// main is the script entry point.
func main() {
	flag.Parse()
	if err := cleanup(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
