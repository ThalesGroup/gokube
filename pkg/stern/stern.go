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
	pb "github.com/cheggaaa/pb/v3"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gemalto/gokube/pkg/download"
)

const (
	DEFAULT_URL           = "https://github.com/stern/stern/releases/download/v%s/stern_%s_windows_amd64.tar.gz"
	LOCAL_EXECUTABLE_NAME = "stern.exe"
)

func bin() string {
	return utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
}

// Version ...
func Version() error {
	fmt.Println("stern version: ")
	cmd := exec.Command(bin(), "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(sternURL string, sternVersion string, bar *pb.ProgressBar) error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if download.IsCurrentVersion(localFile, sternVersion) {
		bar.SetTemplateString(`{{ green "stern" }} ` + sternVersion + ` already up to date`)
		bar.SetTotal(1)
		bar.SetCurrent(1)
		bar.Finish()
		return nil
	}
	_ = os.RemoveAll(localFile)
	_ = os.RemoveAll(download.VersionFile(localFile))
	fileMap := &download.FileMap{Src: LOCAL_EXECUTABLE_NAME, Dst: LOCAL_EXECUTABLE_NAME}
	if _, err := download.FromUrl(sternURL, sternVersion, "stern", []*download.FileMap{fileMap}, filepath.Dir(localFile), bar); err != nil {
		return err
	}
	return download.WriteVersion(localFile, sternVersion)
}

// DeleteExecutable ...
func DeleteExecutable() error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	_ = os.RemoveAll(download.VersionFile(localFile))
	return os.RemoveAll(localFile)
}
