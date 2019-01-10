/*
Copyright 2018 The Knative Authors

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

// gcs.go defines functions to use GCS

package gcs

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"path"
	"os"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"google.golang.org/api/iterator"
)

const (
	latest     = "/latest-build.txt"
)

var client *storage.Client

// Authenticate explicitly sets up authentication for the rest of run
func Authenticate(ctx context.Context, sa string) error {
	var err error
	client, err = storage.NewClient(ctx, option.WithCredentialsFile(sa))
	return err
}

// GetLatestBuildNumber gets the latest build number for the specified log directory
func GetLatestBuildNumber(ctx context.Context, bucketName, logDir string) (int, error) {
	logFilePath := logDir + latest
	log.Printf("Using %s to get latest build number", logFilePath)
	contents, err := ReadGcsFile(ctx, bucketName, logFilePath)
	if err != nil {
		return 0, err
	}
	latestBuild, err := strconv.Atoi(strings.TrimSuffix(string(contents), "\n"))
	if err != nil {
		return 0, err
	}

	return latestBuild, nil
}

//ReadGcsFile reads the specified file using the provided service account
func ReadGcsFile(ctx context.Context, bucketName, filePath string) ([]byte, error) {
	// Create a new GCS client
	o := createStorageObject(bucketName, filePath)
	if _, err := o.Attrs(ctx); err != nil {
		return []byte(fmt.Sprintf("Cannot get attributes of '%s'", filePath)), err
	}
	f, err := o.NewReader(ctx)
	if err != nil {
		return []byte(fmt.Sprintf("Cannot open '%s'", filePath)), err
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte(fmt.Sprintf("Cannot read '%s'", filePath)), err
	}
	return contents, nil
}

// ParseLog parses the log and returns the lines where the checkLog func does not return an empty slice.
// checkLog function should take in the log statement and return a part from that statement that should be in the log output.
func ParseLog(ctx context.Context, bucketName, filePath string, checkLog func(s []string) *string) []string {
	var logs []string

	log.Printf("Parsing '%s'", filePath)
	o := createStorageObject(bucketName, filePath)
	if _, err := o.Attrs(ctx); err != nil {
		log.Printf("Cannot get attributes of '%s', assuming not ready yet: %v", filePath, err)
		return nil
	}
	f, err := o.NewReader(ctx)
	if err != nil {
		log.Fatalf("Error opening '%s': %v", filePath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if s := checkLog(strings.Fields(scanner.Text())); s != nil {
			logs = append(logs, *s)
		}
	}
	return logs
}

// Exist checks if path exist under gcs bucket
func Exist(ctx context.Context, bucketName, filePath string) bool {
	handle := createStorageObject(bucketName, filePath)
	if _, err := handle.Attrs(ctx); nil != err {
		return false
	}
	return true
}

// list child under prefix, use delim to eliminate some files.
// see https://godoc.org/cloud.google.com/go/storage#Query
func List(ctx context.Context, bucketName, prefix, delim string) []string {
	var dirs []string
	objsAttrs := getObjectsAttrs(ctx, bucketName, prefix, delim)
	for _, attrs := range objsAttrs {
		dirs = append(dirs, path.Join(attrs.Prefix, attrs.Name))
	}
	return dirs
}

// list direct child paths(including files and directories)
// Note: to avoid including unnecessary directories, prefix must end with "/".
//  e.g. if there are 2 top directories "foo" and "foobar",
//		then given prefix "foo" with list both from "foo" and "foobar"
func ListDirectChildren(ctx context.Context, bucketName, prefix string) []string {
	return List(ctx, bucketName, prefix, "/")
}

// Copy file from within gcs
func Copy(ctx context.Context, srcBucketName, srcPath, dstBucketName, dstPath string) error {
	src := createStorageObject(srcBucketName, srcPath)
	dst := createStorageObject(dstBucketName, dstPath)

	_, err := dst.CopierFrom(src).Run(ctx)
	return err
}

// Download file from gcs
func Download(ctx context.Context, bucketName, srcPath, dstPath string) error {
	handle := createStorageObject(bucketName, srcPath)
	if _, err := handle.Attrs(ctx); nil != err {
		return err
	}

	dst, _ := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0755)
	// dst, _ := os.OpenFile(dstFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	src, err := handle.NewReader(ctx)
	if err != nil {
		log.Println(2)
		return err
	}
	defer src.Close()
	if _, err = io.Copy(dst, src); nil != err {
		log.Println(3)
		return err
	}

	return nil
}

// Upload file to gcs
func Upload(ctx context.Context, bucketName, dstPath, srcPath string) error {
	src, err := os.Open(srcPath)
	if nil != err {
		return err
	}
	dst := createStorageObject(bucketName, dstPath).NewWriter(ctx)
	defer dst.Close()

	if _, err = io.Copy(dst, src); nil != err {
		return err
	}

	return nil
}

/* Private functions */

// create storage object handle, this step doesn't access internet
func createStorageObject(bucketName, filePath string) *storage.ObjectHandle {
	return client.Bucket(bucketName).Object(filePath)
}

// Query items under given gcs prefix, use delim to eliminate some files.
// see https://godoc.org/cloud.google.com/go/storage#Query
func getObjectsAttrs(ctx context.Context, bucketName, prefix, delim string) []*storage.ObjectAttrs {
	var allAttrs []*storage.ObjectAttrs
	bucketHandle := client.Bucket(bucketName)
	it := bucketHandle.Objects(ctx, &storage.Query{
		Prefix 		:	prefix,
		Delimiter	:	delim,
	})

	for {
		// fmt.Printf("iterating")
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating: %v", err)
		}

		allAttrs = append(allAttrs, attrs)
	}
	
	return allAttrs
}
