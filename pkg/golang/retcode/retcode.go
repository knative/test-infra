package retcode

/*
Copyright 2022 The Knative Authors

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

import "hash/crc32"

// Calc will calculate an POSIX retcode from an error.
func Calc(err error) int {
	if err == nil {
		return 0
	}
	return int(crc32.ChecksumIEEE([]byte(err.Error())))%254 + 1
}
