package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Yaml parsing helpers.

// read template yaml file content
func readTemplate(fp string) string {
	if _, ok := templatesCache[fp]; !ok {
		// get the directory of the currently running file
		_, f, _, _ := runtime.Caller(0)
		content, err := ioutil.ReadFile(path.Join(path.Dir(f), templateDir, fp))
		if err != nil {
			log.Fatalf("Failed read file '%s': '%v'", fp, err)
		}
		templatesCache[fp] = string(content)
	}
	return templatesCache[fp]
}

// getString casts the given interface (expected string) as string.
// An array of length 1 is also considered a single string.
func getString(s interface{}) string {
	if _, ok := s.([]interface{}); ok {
		values := getStringArray(s)
		if len(values) == 1 {
			return values[0]
		}
		log.Fatalf("Entry %v is not a string or string array of size 1", s)
	}
	if str, ok := s.(string); ok {
		return str
	}
	log.Fatalf("Entry %v is not a string", s)
	return ""
}

// getInt casts the given interface (expected int) as int.
func getInt(s interface{}) int {
	if value, ok := s.(int); ok {
		return value
	}
	log.Fatalf("Entry %v is not an integer", s)
	return 0
}

// getBool casts the given interface (expected bool) as bool.
func getBool(s interface{}) bool {
	if value, ok := s.(bool); ok {
		return value
	}
	log.Fatalf("Entry %v is not a boolean", s)
	return false
}

// getInterfaceArray casts the given interface (expected interface array) as interface array.
func getInterfaceArray(s interface{}) []interface{} {
	if interfaceArray, ok := s.([]interface{}); ok {
		return interfaceArray
	}
	log.Fatalf("Entry %v is not an interface array", s)
	return nil
}

// getStringArray casts the given interface (expected string array) as string array.
func getStringArray(s interface{}) []string {
	interfaceArray := getInterfaceArray(s)
	strArray := make([]string, len(interfaceArray))
	for i := range interfaceArray {
		strArray[i] = getString(interfaceArray[i])
	}
	return strArray
}

// getMapSlice casts the given interface (expected MapSlice) as MapSlice.
func getMapSlice(m interface{}) yaml.MapSlice {
	if mm, ok := m.(yaml.MapSlice); ok {
		return mm
	}
	log.Fatalf("Entry %v is not a yaml.MapSlice", m)
	return nil
}

// appendIfUnique appends an element to an array of strings, unless it's already present.
func appendIfUnique(a1 []string, e2 string) []string {
	var res []string
	res = append(res, a1...)
	for _, e1 := range a1 {
		if e1 == e2 {
			return res
		}
	}
	return append(res, e2)
}

func combineSlices(a1 []string, a2 []string) []string {
	var res []string
	res = append(res, a1...)
	for _, e2 := range a2 {
		res = appendIfUnique(res, e2)
	}
	return res
}

// intersectSlices returns intersect of 2 slices
func intersectSlices(a1, a2 []string) []string {
	var res []string
	s1 := sets.NewString(a1...)
	for _, e2 := range a2 {
		if s1.Has(e2) {
			res = append(res, e2)
		}
	}
	return res
}

// exclusiveSlices returns elements in a1 but not in a2
func exclusiveSlices(a1, a2 []string) []string {
	var res []string
	s2 := sets.NewString(a2...)
	for _, e1 := range a1 {
		if !s2.Has(e1) {
			res = append(res, e1)
		}
	}
	return res
}

// getGo112ID returns image identifier for go113 images
func getGo112ID() string {
	return "-go112"
}

// getGo114ID returns image identifier for go114 images
func getGo114ID() string {
	return "-go114"
}

// Get go113 image name from base image name, following the contract of
// [IMAGE]:[DIGEST]-> [IMAGE]-go112:[DIGEST]
func getGo113ImageName(name string) string {
	return stripSuffixFromImageName(name, []string{getGo112ID()})
}

// strip out all suffixes from the image name
func stripSuffixFromImageName(name string, suffixes []string) string {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) != 2 {
		log.Fatalf("image name should contain ':': %q", name)
	}
	for _, s := range suffixes {
		if strings.HasSuffix(parts[0], s) {
			parts[0] = strings.TrimSuffix(parts[0], s)
		}
	}
	return strings.Join(parts, ":")
}

// add suffix to the image name
// e.g. if suffix = "-go112", then [IMAGE]:[DIGEST]-> [IMAGE]-go112:[DIGEST]
func addSuffixToImageName(name string, suffix string) string {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) != 2 {
		log.Fatalf("image name should contain ':': %q", name)
	}
	if !strings.HasSuffix(parts[0], suffix) {
		parts[0] = fmt.Sprintf("%s%s", parts[0], suffix)
	}
	return strings.Join(parts, ":")
}

func getGo112ImageName(name string) string {
	return addSuffixToImageName(name, getGo112ID())
}

func getGo114ImageName(name string) string {
	return addSuffixToImageName(stripSuffixFromImageName(name, []string{getGo112ID()}), getGo114ID())
}

// Consolidate whitelisted and skipped branches with newly added
// whitelisted/skipped. To make the logic easier to maintain, this function
// makes the assumption that the outcome follows these rules:
//   - Special branch logics always apply on master and future branches
// Based on the previous rule, if there is a special branch logic, the 2 Prow
// jobs that serves different branches become:
//   - Standard job definition:
//		- whitelisted: [release-0.1]
//		- skipped: []
//   - Standard job definition + branch special logic #1:
//		- whitelisted: []
//		- skipped: [release-0.1]
// And when there is a new special logic comes up with different list of release
// branches to exclude, for example [release-0.1, release-0.2], then the desired
// outcome becomes:
//   - Standard job definition:
//		- whitelisted: [release-0.1]
//		- skipped: []
//   - Standard job definition + branch special logic #1: (This will never run)
//		- whitelisted: []
//		- skipped: []
//   - Standard job definition + branch special logic #2:
//		- whitelisted: [release-0.2]
//		- skipped: []
//   - Standard job definition + branch special logic #1 + branch special logic #2:
//		- whitelisted: []
//		- skipped: [release-0.1, release-0.2]
// Noted that only jobs with all special branch logics have something in
// skipped, while all other jobs only have whitelisted. This rule also applies
// when there is a third branch specific logic and so on.
// This function takes the logic above, and determines whether generate
// whitelisted or skipped as output.
func consolidateBranches(whitelisted []string, skipped []string, newWhitelisted []string, newSkipped []string) ([]string, []string) {
	var combinedWhitelisted, combinedSkipped []string

	// Do the legacy part(old branches):
	if len(newWhitelisted) > 0 {
		if len(skipped) > 0 {
			// - if previous is skipped(latest), then minus the skipped from current
			// branches, as we want to run exclusive on branches supported currently
			combinedWhitelisted = exclusiveSlices(newWhitelisted, skipped)
		} else if len(whitelisted) > 0 {
			// - if previous is include, then find their intersections, as these are the
			// real supported branches
			combinedWhitelisted = intersectSlices(newWhitelisted, whitelisted)
		} else {
			combinedWhitelisted = newWhitelisted
		}
	} else if len(newSkipped) > 0 { // Then do the pos part(latest)
		if len(skipped) > 0 {
			// - if previous is skipped(latest), then find the combination, as we want to
			// skip all non-supported
			combinedSkipped = combineSlices(newSkipped, skipped)
		} else if len(whitelisted) > 0 {
			// - if previous is include, then minus current branches from included
			combinedWhitelisted = exclusiveSlices(whitelisted, newSkipped)
		} else {
			combinedSkipped = newSkipped
		}
	}
	return combinedWhitelisted, combinedSkipped
}

// recursiveSBL recursively going through specialBranchLogic, and generate job
// at last. Use `i` to keeps track of current index in sbs to be used
func recursiveSBL(repoName string, data interface{}, generateOneJob func(data interface{}), sbs []specialBranchLogic, i int) {
	// Base case, all special branch logics have been applied
	if i == len(sbs) {
		// If there is no branch left, this job shouldn't be generated at all
		if len(getBase(data).Branches) > 0 || len(getBase(data).SkipBranches) > 0 {
			generateOneJob(data)
		}
		return
	}

	sb := sbs[i]
	base := getBase(data)

	origBranches, origSkipBranches := base.Branches, base.SkipBranches
	// Do legacy branches first
	base.Branches, base.SkipBranches = consolidateBranches(origBranches, origSkipBranches, sb.branches, []string{})
	recursiveSBL(repoName, data, generateOneJob, sbs, i+1)
	// Then do latest branches
	base.Branches, base.SkipBranches = consolidateBranches(origBranches, origSkipBranches, []string{}, sb.branches)
	sb.opsNew(base)
	recursiveSBL(repoName, data, generateOneJob, sbs, i+1)
	sb.restore(base)
}

// Multi-value flag parser.

func (a *stringArrayFlag) String() string {
	return strings.Join(*a, ", ")
}

func (a *stringArrayFlag) Set(value string) error {
	*a = append(*a, value)
	return nil
}

// Template helpers.

// gitHubRepo returns the correct reference for the GitHub repository.
func gitHubRepo(data baseProwJobTemplateData) string {
	if repositoryOverride != "" {
		return repositoryOverride
	}
	s := data.RepoURI
	if data.RepoBranch != "" {
		s += "=" + data.RepoBranch
	}
	return s
}

// isNum checks if the given string is a valid number
func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// quote returns the given string quoted if it's not a number, or not a key/value pair, or already quoted.
func quote(s string) string {
	if isNum(s) {
		return s
	}
	if strings.HasPrefix(s, "'") || strings.HasPrefix(s, "\"") || strings.Contains(s, ": ") || strings.HasSuffix(s, ":") {
		return s
	}
	return "\"" + s + "\""
}

// indentBase is a helper function which returns the given array indented.
func indentBase(indentation int, prefix string, indentFirstLine bool, array []string) string {
	s := ""
	if len(array) == 0 {
		return s
	}
	indent := strings.Repeat(" ", indentation)
	for i := 0; i < len(array); i++ {
		if i > 0 || indentFirstLine {
			s += indent
		}
		s += prefix + quote(array[i]) + "\n"
	}
	return s
}

// indentArray returns the given array indented, prefixed by "-".
func indentArray(indentation int, array []string) string {
	return indentBase(indentation, "- ", false, array)
}

// indentKeys returns the given array of key/value pairs indented.
func indentKeys(indentation int, array []string) string {
	return indentBase(indentation, "", false, array)
}

// indentSectionBase is a helper function which returns the given array of key/value pairs indented inside a section.
func indentSectionBase(indentation int, title string, prefix string, array []string) string {
	keys := indentBase(indentation, prefix, true, array)
	if keys == "" {
		return keys
	}
	return title + ":\n" + keys
}

// indentArraySection returns the given array indented inside a section.
func indentArraySection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "- ", array)
}

// indentSection returns the given array of key/value pairs indented inside a section.
func indentSection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "", array)
}

// indentMap returns the given map indented, with each key/value separated by ": "
func indentMap(indentation int, mp map[string]string) string {
	// Extract map keys to keep order consistent.
	keys := make([]string, 0, len(mp))
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	arr := make([]string, len(mp))
	for i := 0; i < len(mp); i++ {
		arr[i] = keys[i] + ": " + quote(mp[keys[i]])
	}
	return indentBase(indentation, "", false, arr)
}

// outputConfig outputs the given line, if not empty, to stdout.
func outputConfig(line string) {
	if strings.TrimSpace(line) != "" {
		fmt.Fprintln(output, strings.TrimRight(line, " "))
		emittedOutput = true
	}
}

// strExists checks if the given string exists in the array
func strExists(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}

/* misc but might be necessary */

// addProjAndRepoIfNeed adds the project and repo if they are new in the metaData map, then return the jobDetailMap
func addProjAndRepoIfNeed(projName string, repoName string) map[string][]string {
	// add project in the metaData
	if _, exists := metaData[projName]; !exists {
		metaData[projName] = make(map[string][]string)
		if !strExists(projNames, projName) {
			projNames = append(projNames, projName)
		}
	}

	// add repo in the project
	jobDetailMap := metaData[projName]
	if _, exists := jobDetailMap[repoName]; !exists {
		if !strExists(repoNames, repoName) {
			repoNames = append(repoNames, repoName)
		}
		jobDetailMap[repoName] = make([]string, 0)
	}
	return jobDetailMap
}

// updateTestCoverageJobDataIfNeeded adds test-coverage job data for the repo if it has go coverage check
func updateTestCoverageJobDataIfNeeded(jobDetailMap *map[string][]string, repoName string) {
	if goCoverageMap[repoName] {
		newJobTypes := append((*jobDetailMap)[repoName], "test-coverage")
		(*jobDetailMap)[repoName] = newJobTypes
		// delete this repoName from the goCoverageMap to avoid it being processed again when we
		// call the function addRemainingTestCoverageJobs
		delete(goCoverageMap, repoName)
	}
}

// addRemainingTestCoverageJobs adds test-coverage jobs data for the repos that haven't been processed.
func addRemainingTestCoverageJobs() {
	// handle repos that only have go coverage
	for repoName, hasGoCoverage := range goCoverageMap {
		if hasGoCoverage {
			jobDetailMap := addProjAndRepoIfNeed(projNames[0], repoName)
			jobDetailMap[repoName] = []string{"test-coverage"}
		}
	}
}

// buildProjRepoStr builds the projRepoStr used in the config file with projName and repoName
func buildProjRepoStr(projName string, repoName string) string {
	projVersion := ""
	if strings.Contains(projName, "-") {
		projNameAndVersion := strings.Split(projName, "-")
		projName = projNameAndVersion[0]
		projVersion = projNameAndVersion[1]
	}
	projRepoStr := repoName
	if projVersion != "" {
		projRepoStr += ("-" + projVersion)
	}
	projRepoStr = projName + "-" + projRepoStr
	return strings.ToLower(projRepoStr)
}

// isReleased returns true for project name that has version
func isReleased(projName string) bool {
	return projNameRegex.FindString(projName) != ""
}

// setOutput set the given file as the output target, then all the output will be written to this file
func setOutput(fileName string) {
	output = os.Stdout
	if fileName == "" {
		return
	}
	configFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Cannot create the configuration file %q: %v", fileName, err)
	}
	configFile.Truncate(0)
	configFile.Seek(0, 0)
	output = configFile
}
