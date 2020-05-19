package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	goVersionMatcher *regexp.Regexp
)

func init() {
	goVersionMatcher = regexp.MustCompile(`go(\d+)[.](\d+)`)
}

// jobDetailMap, key is the repo name, value is the list of job types, like continuous, nightly, etc., as well as custome names
type JobDetailMap map[string][]string

// testGridMetaData saves the meta data needed to generate the final config file.
// key is the main project version, value is another map containing job details
type TestGridMetaData struct {
	md map[string]JobDetailMap
	// projNames save the proj names in a list when parsing the simple config file, for the purpose of maintaining the output sequence
	projNames []string
	// repoNames save the repo names in a list when parsing the simple config file, for the purpose of maintaining the output sequence
	repoNames  []string
	nonAligned []NonAlignedTestGroup
}

type NonAlignedTestGroup struct {
	// DashboardGroup: The things shown at http://testgrid.knative.dev before you hover over anything
	DashboardGroup string
	// DashboardName: This is the thing with multiple tabs/test-groups/whatever-you-call-them
	DashboardName string
	// HumanTabName: Each set of test runs, aka test_group, with the name as shown to the human
	HumanTabName string
	// Used to find the logs
	CIJobName string
	// Extra things that show up in yaml in the test_groups section
	Extra map[string]string
}

type GoVersion struct {
	Major int
	Minor int
}

func (j JobDetailMap) Add(repo, jt string) {
	j.EnsureExists(repo)
	j[repo] = append(j[repo], jt)
}

func NewJobDetailMap() JobDetailMap {
	return make(JobDetailMap)
}

// EnsureExists returns true if already existed or false if newly-created
func (j JobDetailMap) EnsureExists(repo string) bool {
	if _, exists := j[repo]; exists == false {
		j[repo] = make([]string, 0)
		return false
	}
	return true
}

func NewTestGridMetaData() TestGridMetaData {
	return TestGridMetaData{
		md:         make(map[string]JobDetailMap),
		projNames:  make([]string, 0),
		repoNames:  make([]string, 0),
		nonAligned: make([]NonAlignedTestGroup, 0),
	}
}

func NewGoVersion(version string) GoVersion {
	matches := goVersionMatcher.FindStringSubmatch(version)
	if matches == nil || len(matches) != 3 {
		panic(fmt.Errorf("bad version string to NewGoVersion: %q", version))
	}
	a, _ := strconv.Atoi(matches[1])
	b, _ := strconv.Atoi(matches[2])
	return GoVersion{a, b}
}

func (v GoVersion) String() string {
	return fmt.Sprintf("go%d.%d", v.Major, v.Minor)
}

func (v GoVersion) Equals(v2 GoVersion) bool {
	return v.Major == v2.Major && v.Minor == v2.Minor
}
