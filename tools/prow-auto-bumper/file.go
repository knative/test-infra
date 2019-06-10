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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func cdToRootDir() error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	d := strings.TrimSpace(string(output))
	log.Printf("Changing working directory to %s...", d)
	return os.Chdir(d)
}

// update all tags in a byte slice
func (pv *PRVersions) updateAllTags(content []byte, imageFilter *regexp.Regexp) ([]byte, string, []string) {
	var msg string
	var errMsgs []string
	indexes := imageFilter.FindAllSubmatchIndex(content, -1)
	// Not finding any images is not an error.
	if indexes == nil {
		return content, msg, errMsgs
	}

	var res string
	lastIndex := 0
	for _, m := range indexes {
		// append from end of last match to end of image part, including ":"
		res += string(content[lastIndex : m[imageImagePart*2+1]+1])
		// image part of a version, i.e. the portion before ":"
		image := string(content[m[imageImagePart*2]:m[imageImagePart*2+1]])
		// tag part of a version, i.e. the portion after ":"
		tag := string(content[m[imageTagPart*2]:m[imageTagPart*2+1]])
		// m[1] is the end index of current match
		lastIndex = m[1]

		iv := pv.getIndex(image, tag)
		if "" != pv.images[image][iv].newVersion {
			res += pv.images[image][iv].newVersion
			msg += fmt.Sprintf("\nImage: %s\nOld Tag: %s\nNew Tag: %s", image, tag, pv.images[image][iv].newVersion)
		} else {
			errMsg := fmt.Sprintf("Cannot find version for image: '%s:%s'.\n", image, tag)
			log.Println(errMsg)
			errMsgs = append(errMsgs, errMsg)
			res += tag
		}
	}
	res += string(content[lastIndex:])

	return []byte(res), msg, errMsgs
}

// updateFile updates a file in place.
func (pv *PRVersions) updateFile(filename string, imageFilter *regexp.Regexp, dryrun bool) ([]string, error) {
	var errMsgs []string
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return errMsgs, fmt.Errorf("failed to read %s: %v", filename, err)
	}

	newContent, msg, errMsgs := pv.updateAllTags(content, imageFilter)
	if err := run(
		fmt.Sprintf("Update file '%s':%s", filename, msg),
		func() error {
			return ioutil.WriteFile(filename, newContent, 0644)
		},
		dryrun,
	); err != nil {
		return errMsgs, fmt.Errorf("failed to write %s: %v", filename, err)
	}
	return errMsgs, nil
}

func (pv *PRVersions) updateAllFiles(fileFilters []*regexp.Regexp, imageFilter *regexp.Regexp,
	dryrun bool) ([]string, error) {
	var errMsgs []string
	if err := cdToRootDir(); err != nil {
		return errMsgs, fmt.Errorf("failed to change to root dir")
	}

	err := filepath.Walk(".", func(filename string, info os.FileInfo, err error) error {
		for _, ff := range fileFilters {
			if ff.Match([]byte(filename)) {
				msgs, err := pv.updateFile(filename, imageFilter, dryrun)
				errMsgs = append(errMsgs, msgs...)
				if err != nil {
					return fmt.Errorf("Failed to update path %s '%v'", filename, err)
				}
			}
		}
		return nil
	})
	return errMsgs, err
}
