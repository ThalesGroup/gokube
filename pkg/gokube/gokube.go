/*
(c) Copyright 2018, Gemalto. All rights reserved.

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

package gokube

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// GetBinDir ...
func GetBinDir() string {
	path, err := exec.LookPath("gokube")
	if err != nil {
		panic(err)
	}
	if path == "gokube.exe" {
		path = whereAmI()
	} else {
		path = strings.TrimSuffix(path, "\\gokube.exe")
	}
	return path
}

// GetTempDir ...
func GetTempDir() string {
	return GetBinDir() + "/../tmp"
}

// WhereAmI returns a string containing the file name, function name
// and the line number of a specified entry on the call stack
func whereAmI() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
