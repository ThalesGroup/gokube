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

package helm

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/download"
	"os"
	"os/exec"

	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	URL = "https://get.helm.sh/helm-%s-windows-amd64.zip"
)

// Upgrade ...
func Upgrade(chart string, version string, release string, namespace string, configuration string, valuesFile string) {
	var args = []string{"upgrade", "--install", "--devel", release, chart}
	if len(version) > 0 {
		args = append(args, "--version", version)
	}
	if len(namespace) > 0 {
		args = append(args, "--namespace", namespace)
	}
	if len(configuration) > 0 {
		args = append(args, "--set", configuration)
	}
	if len(valuesFile) > 0 {
		args = append(args, "-f", valuesFile)
	}
	fmt.Println("Starting " + chart + " components...")
	cmd := exec.Command("helm", args...)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Delete ...
func Delete(release string) {
	fmt.Println("Deleting " + release + " components...")
	cmd := exec.Command("helm", "delete", release)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// List ...
func List() {
	cmd := exec.Command("helm", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// RepoAdd ...
func RepoAdd(name string, repo string) {
	cmd := exec.Command("helm", "repo", "add", name, repo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//RepoUpdate ...
func RepoUpdate() {
	cmd := exec.Command("helm", "repo", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// RepoRemove ...
func RepoRemove(name string) {
	cmd := exec.Command("helm", "repo", "remove", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//Version ...
func Version() {
	fmt.Print("helm version: ")
	cmd := exec.Command("helm", "version", "--client", "--short")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//Version ...
func PluginsVersion() {
	fmt.Print("helm plugins version:\n")
	cmd := exec.Command("helm", "plugin", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(dst string, helmVersion string) {
	if _, err := os.Stat(gokube.GetBinDir() + "/helm.exe"); os.IsNotExist(err) {
		download.FromUrl("helm "+helmVersion, URL, helmVersion)
		utils.MoveFile(gokube.GetTempDir()+"/windows-amd64/helm.exe", dst+"/helm.exe")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeleteExecutable ...
func DeleteExecutable() {
	utils.RemoveFile(gokube.GetBinDir() + "/helm.exe")
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() {
	// This directory contains helm plugins and repo definitions and caches
	utils.RemoveDir(utils.GetAppDataHome() + "/helm")
}

// ResetWorkingDirectory ...
func ResetWorkingDirectory() {
	utils.RemoveFile(utils.GetAppDataHome() + "/helm/repositories.yaml")
	utils.RemoveFile(utils.GetAppDataHome() + "/helm/repositories.lock")
}
