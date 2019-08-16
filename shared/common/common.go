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

package common

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const allUsersFullPermission = 0777

// CreateDir creates dir if does not exist.
// The created dir will have the permission bits as 0777, which means everyone can read/write/execute it.
func CreateDir(dirPath string) error {
	return CreateDirWithFileMode(dirPath, allUsersFullPermission)
}

// CreateDirWithFileMode creates dir if does not exist.
// The created dir will have the permission bits as perm, which is the standard Unix rwxrwxrwx permissions.
func CreateDirWithFileMode(dirPath string, perm os.FileMode) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, perm); err != nil {
			return fmt.Errorf("Failed to create directory: %v", err)
		}
	}
	return nil
}

// GetRootDir gets directory of git root
func GetRootDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// CDToRootDir change directory to git root dir
func CDToRootDir() error {
	d, err := GetRootDir()
	if nil != err {
		return err
	}
	return os.Chdir(d)
}

// GetString casts the given interface (expected string) as string.
// An array of length 1 is also considered a single string.
func GetString(s interface{}) string {
	if _, ok := s.([]interface{}); ok {
		values := GetStringArray(s)
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

// GetInt casts the given interface (expected int) as int.
func GetInt(s interface{}) int {
	if value, ok := s.(int); ok {
		return value
	}
	log.Fatalf("Entry %v is not an integer", s)
	return 0
}

// GetBool casts the given interface (expected bool) as bool.
func GetBool(s interface{}) bool {
	if value, ok := s.(bool); ok {
		return value
	}
	log.Fatalf("Entry %v is not a boolean", s)
	return false
}

// GetInterfaceArray casts the given interface (expected interface array) as interface array.
func GetInterfaceArray(s interface{}) []interface{} {
	if interfaceArray, ok := s.([]interface{}); ok {
		return interfaceArray
	}
	log.Fatalf("Entry %v is not an interface array", s)
	return nil
}

// GetStringArray casts the given interface (expected string array) as string array.
func GetStringArray(s interface{}) []string {
	interfaceArray := GetInterfaceArray(s)
	strArray := make([]string, len(interfaceArray))
	for i := range interfaceArray {
		strArray[i] = GetString(interfaceArray[i])
	}
	return strArray
}

// GetMapSlice casts the given interface (expected MapSlice) as MapSlice.
func GetMapSlice(m interface{}) yaml.MapSlice {
	if mm, ok := m.(yaml.MapSlice); ok {
		return mm
	}
	log.Fatalf("Entry %v is not a yaml.MapSlice", m)
	return nil
}

// IsNum checks if the given string is a valid number
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Quote returns the given string quoted if it's not a number, or not a key/value pair, or already quoted.
func Quote(s string) string {
	if IsNum(s) {
		return s
	}
	if strings.HasPrefix(s, "'") || strings.HasPrefix(s, "\"") || strings.Contains(s, ": ") || strings.HasSuffix(s, ":") {
		return s
	}
	return "\"" + s + "\""
}

// StrExists checks if the given string exists in the array
func StrExists(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}
