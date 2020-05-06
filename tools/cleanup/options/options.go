package options

import (
	"flag"
	"strings"
)

type Options struct {
	ProjectResourceYaml  strSliceArg
	Project              strSliceArg
	ReProjectName        string
	DaysToKeepImages     int
	HoursToKeepClusters  int
	Registry             string
	ServiceAccount       string
	ConcurrentOperations int
	DryRun               bool
}

type strSliceArg []string

func (ss *strSliceArg) String() string {
	return strings.Join(*ss, ", ")
}

func (ss *strSliceArg) Set(val string) error {
	*ss = append(*ss, val)
	return nil
}

func (o *Options) AddOptions() {
	flag.Var(&o.ProjectResourceYaml, "project-resource-yaml", "Resources file containing the names of the projects to be cleaned up.")
	flag.Var(&o.Project, "project", "Project to be cleaned up.")
	flag.StringVar(&o.ReProjectName, "re-project-name", "knative-boskos-[a-zA-Z0-9]+", "Regular expression for filtering project names from the resources file.")
	flag.IntVar(&o.DaysToKeepImages, "days-to-keep-images", 365, "Images older than this amount of days will be deleted (defaults to 1 year, -1 means 'forever').")
	flag.IntVar(&o.HoursToKeepClusters, "hours-to-keep-clusters", 720, "Clusters older than this amount of hours will be deleted (defaults to 1 month, -1 means 'forever').")
	flag.StringVar(&o.Registry, "gcr", "gcr.io", "The registry hostname to use (defaults to gcr.io; currently only GCR is supported).")
	flag.StringVar(&o.ServiceAccount, "service-account", "", "Specify the key file of the service account to use.")
	flag.IntVar(&o.ConcurrentOperations, "concurrent-operations", 10, "How many deletion operations to run concurrently (defaults to 10).")
	flag.BoolVar(&o.DryRun, "dry-run", false, "Performs a dry run for all deletion functions.")
}
