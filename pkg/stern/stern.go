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

package stern

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	URL = "https://github.com/wercker/stern/releases/download/%s/stern_windows_amd64.exe"
)

// Version ...
func Version() {
	fmt.Println("stern version: ")
	cmd := exec.Command("stern", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(dst string, sternVersion string) {
	if _, err := os.Stat(gokube.GetBinDir() + "/stern.exe"); os.IsNotExist(err) {
		download.DownloadFromUrl("stern v"+sternVersion, URL, sternVersion)
		utils.MoveFile(gokube.GetTempDir()+"/stern_windows_amd64.exe", dst+"/stern.exe")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeleteExecutable ...
func DeleteExecutable() {
	utils.RemoveFile(gokube.GetBinDir() + "/stern.exe")
}
