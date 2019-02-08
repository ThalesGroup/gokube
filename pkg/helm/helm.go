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
	"os"
	"os/exec"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	URL = "https://storage.googleapis.com/kubernetes-helm/helm-%s-windows-amd64.tar.gz"
)

// Upgrade ...
func Upgrade(chart string, release string) {
	fmt.Println("Starting " + chart + " components...")
	cmd := exec.Command("helm", "upgrade", "--install", "--devel", release, chart)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Delete ...
func Delete(release string) {
	fmt.Println("Deleting " + release + " components...")
	cmd := exec.Command("helm", "delete", release, "--purge")
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// UpgradeWithNamespaceVersionAndConfiguration ...
func UpgradeWithNamespaceVersionAndConfiguration(name string, namespace string, version string, configuration string, chart string) {
	fmt.Println("Starting " + chart + " components...")
	cmd := exec.Command("helm", "upgrade", "--install", "--devel", name, "--namespace", namespace, "--version", version, "--set", configuration, chart)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// UpgradeWithConfiguration ...
func UpgradeWithConfiguration(name string, namespace string, configuration string, chart string, version string) {
	fmt.Println("Starting " + chart + " components...")
	cmd := exec.Command("helm", "install", chart, "--name", name, "--namespace", namespace, "--set", configuration, "--version", version)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// UpgradeWithValues ...
func UpgradeWithValues(namespace string, values string, chart string) {
	fmt.Println("Starting " + chart + " components...")
	cmd := exec.Command("helm", "install", chart, "--namespace", namespace, "-f", values)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// UpgradeWithNamespaceAndVersion ...
func UpgradeWithNamespaceAndVersion(name string, namespace string, version string, chart string) {
	fmt.Println("Starting " + chart + " components...")
	cmd := exec.Command("helm", "upgrade", "--install", "--devel", name, "--namespace", namespace, "--version", version, chart)
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Init ...
func Init() {
	cmd := exec.Command("helm", "init", "--upgrade", "--wait", "--tiller-connection-timeout", "600")
	cmd.Stdout = os.Stdout
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

// DownloadExecutable ...
func DownloadExecutable(dst string, helmVersion string) {
	if _, err := os.Stat(gokube.GetBinDir() + "/helm.exe"); os.IsNotExist(err) {
		download.DownloadFromUrl("helm "+helmVersion, URL, helmVersion)
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
	utils.CleanDir(utils.GetUserHome() + "/.helm")
}
