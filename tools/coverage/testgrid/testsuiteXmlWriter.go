/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package testgrid

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/knative/test-infra/tools/coverage/artifacts"
	"github.com/knative/test-infra/tools/coverage/calc"
	"github.com/knative/test-infra/tools/coverage/logUtil"
)

type Property struct {
	XMLName string `xml:"property"`
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
}

type Properties struct {
	XMLName      string `xml:"properties"`
	PropertyList []Property
}

type TestCase struct {
	XMLName      string     `xml:"testcase"`
	ClassName    string     `xml:"class_name,attr"`
	Name         string     `xml:"name,attr"`
	Time         string     `xml:"time,attr"`
	Failure      bool       `xml:"failure,omitempty"`
	PropertyList Properties `xml:"properties"`
}

// NewTestCase constructs the TestCase struct
func NewTestCase(targetName, coverage string, failure bool) *TestCase {
	properties := &Properties{}
	properties.PropertyList = append(properties.PropertyList, Property{"", "coverage", coverage})

	return &TestCase{"", "go_coverage", targetName, "0", failure, *properties}
}

type Testsuite struct {
	XMLName   string     `xml:"testsuite"`
	Testcases []TestCase `xml:"testsuite"`
}

// addTestCase adds one test case to testsuite
func (ts *Testsuite) addTestCase(tc TestCase) {
	ts.Testcases = append(ts.Testcases, tc)
}

// toTestsuite populates Testsuite struct with data from CoverageList and actual file
// directories from OS
func toTestsuite(g *calc.CoverageList, dirs []string) (ts *Testsuite) {
	ts = &Testsuite{}
	g.Summarize()
	covThresInt := g.CovThresInt()
	ts.addTestCase(*NewTestCase("OVERALL", g.PercentageForTestgrid(),
		g.IsCoverageLow(covThresInt)))

	fmt.Println("")
	log.Println("Constructing Testsuite Struct for Testgrid")
	for _, cov := range *g.Group() {
		coverage := cov.PercentageForTestgrid()
		if coverage != "" {
			ts.addTestCase(*NewTestCase(cov.Name(), coverage, cov.IsCoverageLow(covThresInt)))
		} else {
			log.Printf("Skipping file %s as it has no coverage data.\n", cov.Name())
		}
	}

	for _, dir := range dirs {
		dirCov := g.Subset(dir)
		coverage := dirCov.PercentageForTestgrid()
		if coverage != "" {
			ts.addTestCase(*NewTestCase(dir, coverage, dirCov.IsCoverageLow(covThresInt)))
		} else {
			log.Printf("Skipping directory %s as it has no files with coverage data.\n", dir)
		}
	}
	log.Println("Finished Constructing Testsuite Struct for Testgrid")
	fmt.Println("")
	return
}

// ProfileToTestsuiteXML uses coverage profile (and it's corresponding stdout) to produce junit xml
// which serves as the input for test coverage testgrid
func ProfileToTestsuiteXML(arts *artifacts.LocalArtifacts, covThres int) {
	groupCov := calc.CovList(
		artifacts.NewProfileReader(arts.ProfileReader()),
		nil,
		nil,
		covThres,
	)
	f, err := os.Create(arts.JunitXmlForTestgridPath())
	if err != nil {
		logUtil.LogFatalf("Cannot create file: %v", err)
	}
	defer f.Close()

	ts := toTestsuite(groupCov, groupCov.GetDirs())
	output, err := xml.MarshalIndent(ts, "", "    ")
	if err != nil {
		logUtil.LogFatalf("error: %v\n", err)
	}

	f.Write(output)
}
