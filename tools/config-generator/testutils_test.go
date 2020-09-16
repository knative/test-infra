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

package main

import (
	"bytes"
	"os"
	"testing"
)

var outputBuffer bytes.Buffer

// logFatalCalls tracks the number of logFatalf calls that occurred within a test
var logFatalCalls int

func logFatalfMock(format string, v ...interface{}) {
	logFatalCalls++
}

func ResetOutput() {
	outputBuffer = bytes.Buffer{}
	output = newOutputter(&outputBuffer)
}

func GetOutput() string {
	return outputBuffer.String()
}

func TestMain(m *testing.M) {
	ResetOutput() // Redirect output prior to each test.
	logFatalf = logFatalfMock
	logFatalCalls = 0
	os.Exit(m.Run())
}
