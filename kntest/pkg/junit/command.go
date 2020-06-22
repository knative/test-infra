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
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"knative.dev/test-infra/pkg/junit"
)

type option struct {
	suite  string
	name   string
	errMsg string
	dest   string
}

func addOptions(cmd *cobra.Command, opt *option) {
	pf := cmd.Flags()
	pf.StringVar(&opt.suite, "suite", "", "Name of suite")
	pf.StringVar(&opt.name, "name", "", "Name of test")
	pf.StringVar(&opt.errMsg, "err-msg", "", "Error message, empty means test passed, default empty")
	pf.StringVar(&opt.dest, "dest", "junit_result.xml", "Where junit xml writes to")
}

func AddCommands(topLevel *cobra.Command) {
	opt := &option{}
	var junitCmd = &cobra.Command{
		Use:   "junit",
		Short: "Commands for manipulating junit formatted xml files.",
		Run: func(cmd *cobra.Command, args []string) {
			suites := junit.TestSuites{}
			suite := junit.TestSuite{Name: opt.suite}
			tc := junit.TestCase{Name: opt.name}
			if opt.errMsg != "" {
				// errMsg := html.EscapeString(errMsg)
				tc.Failure = &opt.errMsg
			}
			suite.AddTestCase(tc)
			// Ignore the error as it only happens if the test suite name already exists.
			suites.AddTestSuite(&suite)
			contents, err := suites.ToBytes("", "")
			if err != nil {
				log.Fatal(err)
			}
			if err := ioutil.WriteFile(opt.dest, contents, 0644); err != nil {
				log.Fatalf("Error writing to file %q: %v", opt.dest, err)
			}
		},
	}

	addOptions(junitCmd, opt)
	topLevel.AddCommand(junitCmd)
}
