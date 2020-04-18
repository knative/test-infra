package main

import (
	"fmt"
	"log"
	"strings"
)

// General helpers.

// createCommand returns an array with the command to run and its arguments.
func createCommand(data baseProwJobTemplateData) []string {
	c := []string{data.Command}
	// Prefix the pre-command if present.
	if preCommand != "" {
		c = append([]string{preCommand}, c...)
	}
	return append(c, data.Args...)
}

// addEnvToJob adds the given key/pair environment variable to the job.
func addEnvToJob(data *baseProwJobTemplateData, key, value string) {
	// Value should always be string. Add quotes if we get a number
	if isNum(value) {
		value = "\"" + value + "\""
	}

	(*data).Env = append((*data).Env, []string{"- name: " + key, "  value: " + value}...)
}

// addLabelToJob adds extra labels to a job
func addLabelToJob(data *baseProwJobTemplateData, key, value string) {
	(*data).Labels = append((*data).Labels, []string{key + ": " + value}...)
}

// addPubsubLabelsToJob adds the pubsub labels so the prow job message will be picked up by test-infra monitoring
func addMonitoringPubsubLabelsToJob(data *baseProwJobTemplateData, runID string) {
	addLabelToJob(data, "prow.k8s.io/pubsub.project", "knative-tests")
	addLabelToJob(data, "prow.k8s.io/pubsub.topic", "knative-monitoring")
	addLabelToJob(data, "prow.k8s.io/pubsub.runID", runID)
}

// addVolumeToJob adds the given mount path as volume for the job.
func addVolumeToJob(data *baseProwJobTemplateData, mountPath, name string, isSecret bool, defaultMode string) {
	(*data).VolumeMounts = append((*data).VolumeMounts, []string{"- name: " + name, "  mountPath: " + mountPath}...)
	if isSecret {
		(*data).VolumeMounts = append((*data).VolumeMounts, "  readOnly: true")
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
	(*data).Volumes = append((*data).Volumes, s...)
}

// configureServiceAccountForJob adds the necessary volumes for the service account for the job.
func configureServiceAccountForJob(data *baseProwJobTemplateData) {
	if data.ServiceAccount == "" {
		return
	}
	p := strings.Split(data.ServiceAccount, "/")
	if len(p) != 4 || p[0] != "" || p[1] != "etc" || p[3] != "service-account.json" {
		log.Fatalf("Service account path %q is expected to be \"/etc/<name>/service-account.json\"", data.ServiceAccount)
	}
	name := p[2]
	addVolumeToJob(data, "/etc/"+name, name, true, "")
}

// addExtraEnvVarsToJob adds extra environment variables to a job.
func addExtraEnvVarsToJob(envVars []string, data *baseProwJobTemplateData) {
	for _, env := range envVars {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			log.Fatalf("Environment variable %q is expected to be \"key=value\"", env)
		}
		addEnvToJob(data, pair[0], pair[1])
	}
}

// setupDockerInDockerForJob enables docker-in-docker for the given job.
func setupDockerInDockerForJob(data *baseProwJobTemplateData) {
	addVolumeToJob(data, "/docker-graph", "docker-graph", false, "")
	addEnvToJob(data, "DOCKER_IN_DOCKER_ENABLED", "\"true\"")
	(*data).SecurityContext = []string{"privileged: true"}
}

// setReourcesReqForJob sets resource requirement for job
func setReourcesReqForJob(r resources, data *baseProwJobTemplateData) {
	var rsToStr func(key string, rd resourcesDef) []string
	rsToStr = func(key string, rd resourcesDef) []string {
		var res []string
		if rd.CPU != "" || rd.Memory != "" {
			res = append(res, fmt.Sprintf("  %s:", key))
			res = append(res, fmt.Sprintf("    cpu: %s", rd.CPU))
			res = append(res, fmt.Sprintf("    memory: %s", rd.Memory))
		}
		return res
	}
	data.Resources = append(data.Resources, rsToStr("limits", r.Limits)...)
	data.Resources = append(data.Resources, rsToStr("requests", r.Requests)...)
}
