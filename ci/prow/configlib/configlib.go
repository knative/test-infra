package configlib

import (
	"log"
	"strings"

	"knative.dev/test-infra/shared/common"
)

// BaseProwJobTemplateData contains basic data about a Prow job.
type BaseProwJobTemplateData struct {
	OrgName             string
	RepoName            string
	RepoNameForJob      string
	GcsBucket           string
	GcsLogDir           string
	GcsPresubmitLogDir  string
	RepoURI             string
	RepoBranch          string
	CloneURI            string
	SecurityContext     []string
	SkipBranches        []string
	Branches            []string
	DecorationConfig    []string
	ExtraRefs           []string
	Command             string
	Args                []string
	Env                 []string
	Volumes             []string
	VolumeMounts        []string
	Timeout             int
	AlwaysRun           bool
	LogsDir             string
	PresubmitLogsDir    string
	TestAccount         string
	ServiceAccount      string
	ReleaseGcs          string
	GoCoverageThreshold int
	Image               string
	Year                int
	Labels              []string
	PathAlias           string
	Optional            string
}

// PresubmitJobTemplateData contains data about a presubmit Prow job.
type PresubmitJobTemplateData struct {
	Base                 BaseProwJobTemplateData
	PresubmitJobName     string
	PresubmitPullJobName string
	PresubmitPostJobName string
	PresubmitCommand     []string
}

// PostsubmitJobTemplateData contains data about a postsubmit Prow job.
type PostsubmitJobTemplateData struct {
	Base              BaseProwJobTemplateData
	PostsubmitJobName string
}

// AddEnvToJob adds the given key/pair environment variable to the job.
func AddEnvToJob(d *BaseProwJobTemplateData, key, value string) {
	// Value should always be string. Add quotes if we get a number
	if common.IsNum(value) {
		value = "\"" + value + "\""
	}

	(*d).Env = append((*d).Env, []string{"- name: " + key, "  value: " + value}...)
}

// AddLabelToJob adds extra labels to a job
func AddLabelToJob(d *BaseProwJobTemplateData, key, value string) {
	(*d).Labels = append((*d).Labels, []string{key + ": " + value}...)
}

// AddVolumeToJob adds the given mount path as volume for the job.
func AddVolumeToJob(d *BaseProwJobTemplateData, mountPath, name string, isSecret bool, defaultMode string) {
	(*d).VolumeMounts = append((*d).VolumeMounts, []string{"- name: " + name, "  mountPath: " + mountPath}...)
	if isSecret {
		(*d).VolumeMounts = append((*d).VolumeMounts, "  readOnly: true")
	}
	s := []string{"- name: " + name}
	if isSecret {
		arr := []string{"  secret:", "    secretName: " + name}
		if len(defaultMode) > 0 {
			arr = append(arr, "    defaultMode: "+defaultMode)
		}
		s = append(s, arr...)
	} else {
		s = append(s, "  emptyDir: {}")
	}
	(*d).Volumes = append((*d).Volumes, s...)
}

// ConfigureServiceAccountForJob adds the necessary volumes for the service account for the job.
func ConfigureServiceAccountForJob(d *BaseProwJobTemplateData) {
	if d.ServiceAccount == "" {
		return
	}
	p := strings.Split(d.ServiceAccount, "/")
	if len(p) != 4 || p[0] != "" || p[1] != "etc" || p[3] != "service-account.json" {
		log.Fatalf("Service account path %q is expected to be \"/etc/<name>/service-account.json\"", d.ServiceAccount)
	}
	name := p[2]
	AddVolumeToJob(d, "/etc/"+name, name, true, "")
}

// AddExtraEnvVarsToJob adds extra environment variables to a job.
func AddExtraEnvVarsToJob(d *BaseProwJobTemplateData, envVars []string) {
	for _, env := range envVars {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			log.Fatalf("Environment variable %q is expected to be \"key=value\"", env)
		}
		AddEnvToJob(d, pair[0], pair[1])
	}
}
