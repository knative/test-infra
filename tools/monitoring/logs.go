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

// ErrorLog stores the error pattern and the corresponding error message found in the log
type ErrorLog struct {
	Pattern string
	Msg     string
}

// compilePatterns compiles the patterns from string to Regexp. In addition it returns the list of
// patterns that cannot be compiled
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

// collectMatches collects error messages that matches the patterns from text.
func collectMatches(regexps []regexp.Regexp, text []byte) []ErrorLog {
	var errorLogs []ErrorLog
	for _, r := range regexps {
		found := r.Find(text)
		if found != nil {
			errorLogs = append(errorLogs, ErrorLog{
				Pattern: r.String(),
				Msg:     string(found),
			})
		}
	}
	return errorLogs
}

// ParseLog fetches build log from given url and checks it against given error patterns. Return
// all found error patterns and error messages in pairs.
func ParseLog(url string, patterns []string) ([]ErrorLog, error) {
	content, err := getFileBytes(url)
	if err != nil {
		return nil, err
	}
	regexps, badPatterns := compilePatterns(patterns)
	if len(badPatterns) != 0 {
		log.Printf("The following patterns cannot be compiled: %v", badPatterns)
	}

	return collectMatches(regexps, content), nil
}
