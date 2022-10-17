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
	"path/filepath"

	"github.com/gemalto/gokube/pkg/utils"
)

const (
	DEFAULT_URL           = "https://get.helm.sh/helm-%s-windows-amd64.zip"
	LOCAL_EXECUTABLE_NAME = "helm.exe"
)

// Upgrade ...
func Upgrade(chart string, version string, release string, namespace string, configuration string, valuesFile string) error {
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
	return cmd.Run()
}

// RepoAdd ...
func RepoAdd(name string, repo string) error {
	cmd := exec.Command("helm", "repo", "add", name, repo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RepoUpdate ...
func RepoUpdate() error {
	cmd := exec.Command("helm", "repo", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Version ...
func Version() error {
	fmt.Print("helm version: ")
	cmd := exec.Command("helm", "version", "--client", "--short")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PluginsVersion ...
func PluginsVersion() error {
	fmt.Println("helm plugins version:")
	cmd := exec.Command("helm", "plugin", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(helmURL string, helmVersion string) error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		fileMap := &download.FileMap{Src: "windows-amd64" + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME, Dst: LOCAL_EXECUTABLE_NAME}
		_, err = download.FromUrl(helmURL, helmVersion, "helm", []*download.FileMap{fileMap}, filepath.Dir(localFile))
		if err != nil {
			return nil
		}
	}
	return nil
}

// DeleteExecutable ...
func DeleteExecutable() error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	return os.RemoveAll(localFile)
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() error {
	// This directory contains helm plugins and repo definitions and caches
	return os.RemoveAll(utils.GetAppDataHome() + string(os.PathSeparator) + "helm")
}

// ResetWorkingDirectory ...
func ResetWorkingDirectory() error {
	err := os.RemoveAll(utils.GetAppDataHome() + string(os.PathSeparator) + "helm" + string(os.PathSeparator) + "repositories.yaml")
	if err != nil {
		return err
	}
	err = os.RemoveAll(utils.GetAppDataHome() + string(os.PathSeparator) + "helm" + string(os.PathSeparator) + "repositories.lock")
	if err != nil {
		return err
	}
	return nil
}
