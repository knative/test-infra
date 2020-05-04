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

package junit

import (
	"html"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"knative.dev/pkg/test/junit"
)

var (
	suite  string
	name   string
	errMsg string
	dest   string
)

func addOptions(cmd *cobra.Command) {
	pf := cmd.Flags()
	pf.StringVar(&suite, "suite", "", "Name of suite")
	pf.StringVar(&name, "name", "", "Name of test")
	pf.StringVar(&errMsg, "err-msg", "", "Error message, empty means test passed, default empty")
	pf.StringVar(&dest, "dest", "junit_result.xml", "Where junit xml writes to")
}

func AddCommands(topLevel *cobra.Command) {
	var junitCmd = &cobra.Command{
		Use:   "junit",
		Short: "Commands for manipulating junit formatted xml files.",
		Run: func(cmd *cobra.Command, args []string) {
			suites := junit.TestSuites{}
			suite := junit.TestSuite{Name: suite}
			errMsg := html.EscapeString(errMsg)
			suite.AddTestCase(junit.TestCase{
				Name:    name,
				Failure: &errMsg,
			})
			// Ignore the error as it only happens if the test suite name already exists.
			suites.AddTestSuite(&suite)
			contents, err := suites.ToBytes("", "")
			if err != nil {
				log.Fatal(err)
			}
			if err := ioutil.WriteFile(dest, contents, 0644); err != nil {
				log.Fatalf("Error writing to file %q: %v", dest, err)
			}
		},
	}

	addOptions(junitCmd)
	topLevel.AddCommand(junitCmd)
}
