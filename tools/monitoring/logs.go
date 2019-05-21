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

package main

import (
	"log"
	"regexp"
)

// ErrorInLog stores the error pattern and the corresponding error message found in the log
type ErrorInLog struct {
	ErrorPattern string
	ErrorMsg     string
}

func compilePatterns(patterns []string) ([]regexp.Regexp, []string) {
	var regexps []regexp.Regexp
	var badPatterns []string // patterns that cannot be compiled into regex

	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("Error compiling pattern [%s]: %v", pattern, err)
			badPatterns = append(badPatterns, pattern)
		} else {
			regexps = append(regexps, *r)
		}
	}
	return regexps, badPatterns
}

func findMatches(regexps []regexp.Regexp, text []byte) []ErrorInLog {
	var errorLogs []ErrorInLog
	for _, r := range regexps {
		found := r.Find(text)
		if found != nil {
			errorLogs = append(errorLogs, ErrorInLog{
				ErrorPattern: r.String(),
				ErrorMsg:     string(found),
			})
		}
	}
	return errorLogs
}

// ParseLog fetches build log from given url and checks it against given error patterns. Return
// all found error patterns and error messages in pairs.
func ParseLog(url string, patterns []string) ([]ErrorInLog, error) {
	content, err := getFileBytes(url)
	if err != nil {
		return nil, err
	}
	regexps, badPatterns := compilePatterns(patterns)
	if len(badPatterns) != 0 {
		log.Printf("The following pattern cannot be compiled: %v", badPatterns)
		// TODO: after email feature is done, send an email to oncall about the pattern issues
	}

	return findMatches(regexps, content), nil
}
