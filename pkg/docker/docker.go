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
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	URL     = "https://download.docker.com/win/static/edge/x86_64/docker-%s-ce.zip"
	VERSION = "17.10.0"
)

// LoadImages ...
func LoadImages(imagesDir string, envVars []utils.EnvVar) {
	fileList := []string{}
	err := filepath.Walk(imagesDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error while parsing %s: %v\n", path, err)
		} else {
			if f != nil && !f.IsDir() {
				fileList = append(fileList, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error while parsing %s: %v\n", imagesDir, err)
	}
	for _, file := range fileList {

		LoadImage(file, envVars)
	}
}

// LoadImage ...
func LoadImage(image string, envVars []utils.EnvVar) {
	cmd := exec.Command("docker", "load", "-i", image)
	cmd.Env = append(os.Environ())
	for _, element := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", element.Name, element.Value))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// TagImage ...
func TagImage(image string, tag string, envVars []utils.EnvVar) {
	cmd := exec.Command("docker", "tag", image, tag)
	cmd.Env = append(os.Environ())
	for _, element := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", element.Name, element.Value))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Version ...
func Version() {
	fmt.Println("docker version: ")
	cmd := exec.Command("docker", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Download ...
func Download(dst string) {
	if _, err := os.Stat(gokube.GetBinDir() + "/docker.exe"); os.IsNotExist(err) {
		download.DownloadFromUrl("docker v"+VERSION, URL, VERSION)
		utils.MoveFile(gokube.GetTempDir()+"/docker/docker.exe", dst+"/docker.exe")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// RemoveImage ...
func RemoveImage(image string, envVars []utils.EnvVar) {
	cmd := exec.Command("docker", "rmi", image)
	cmd.Env = append(os.Environ())
	for _, element := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", element.Name, element.Value))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Purge ...
func Purge() {
	utils.RemoveFile(gokube.GetBinDir() + "/docker.exe")
}
