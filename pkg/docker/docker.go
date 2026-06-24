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

package docker

import (
	"fmt"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	DEFAULT_URL           = "https://download.docker.com/win/static/stable/x86_64/docker-%s.zip"
	LOCAL_EXECUTABLE_NAME = "docker.exe"
)

func bin() string {
	return utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
}

// Version ...
func Version() error {
	fmt.Println("docker version:")
	cmd := exec.Command(bin(), "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(dockerURL string, dockerVersion string, bar *pb.ProgressBar) error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if download.IsCurrentVersion(localFile, dockerVersion) {
		bar.SetTemplateString(`{{ green "docker" }} ` + dockerVersion + ` already up to date`)
		bar.SetTotal(1)
		bar.SetCurrent(1)
		bar.Finish()
		return nil
	}
	_ = os.RemoveAll(localFile)
	_ = os.RemoveAll(download.VersionFile(localFile))
	fileMap := &download.FileMap{Src: "docker" + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME, Dst: LOCAL_EXECUTABLE_NAME}
	if _, err := download.FromUrl(dockerURL, dockerVersion, "docker", []*download.FileMap{fileMap}, filepath.Dir(localFile), bar); err != nil {
		return err
	}
	return download.WriteVersion(localFile, dockerVersion)
}

// DeleteExecutable ...
func DeleteExecutable() error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	_ = os.RemoveAll(download.VersionFile(localFile))
	return os.RemoveAll(localFile)
}

// InitWorkingDirectory ...
func InitWorkingDirectory() error {
	var dockerHome = utils.GetUserHome() + string(os.PathSeparator) + ".docker"
	var configJsonPath = dockerHome + string(os.PathSeparator) + "config.json"
	_, err := os.Stat(configJsonPath)
	if err == nil {
		return nil
	}
	err = utils.CreateDirs(dockerHome)
	if err != nil {
		return err
	}
	configFile, err := os.Create(configJsonPath)
	defer utils.CloseFile(configFile)
	if err != nil {
		return err
	}
	_, _ = configFile.WriteString("{}")
	err = configFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() error {
	// Delete and recreate will not work if .docker is a symlink !
	return utils.CleanDir(utils.GetUserHome() + string(os.PathSeparator) + ".docker")
}
