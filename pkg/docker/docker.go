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
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	DEFAULT_URL           = "https://github.com/cvila84/docker-cli-builder/releases/download/%s/docker.exe"
	LOCAL_EXECUTABLE_NAME = "docker.exe"
)

// Version ...
func Version() error {
	fmt.Println("docker version:")
	cmd := exec.Command("docker", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(dockerURL string, dockerVersion string) error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		fileMap := &download.FileMap{Src: LOCAL_EXECUTABLE_NAME, Dst: LOCAL_EXECUTABLE_NAME}
		_, err = download.FromUrl(dockerURL, dockerVersion, "docker", []*download.FileMap{fileMap}, filepath.Dir(localFile))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteExecutable ...
func DeleteExecutable() error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
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
