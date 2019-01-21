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

package kubectl

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	URL = "https://storage.googleapis.com/kubernetes-release/release/%s/bin/windows/amd64/kubectl.exe"
)

// ConfigUseContext ...
func ConfigUseContext(context string) {
	cmd := exec.Command("kubectl", "config", "use-context", context)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Patch ...
func Patch(namespace string, resourceType string, resourceName string, patch string) {
	cmd := exec.Command("kubectl", "--namespace", namespace, "patch", resourceType, resourceName, "-p", patch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Apply ...
func Apply(file string, namespace string) {
	cmd := exec.Command("kubectl", "create", "-f", file, "--namespace", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Version ...
func Version() {
	fmt.Println("kubectl version: ")
	cmd := exec.Command("kubectl", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// DisabledNetworkPolicy ...
func DisabledNetworkPolicy() {
	cmd := exec.Command("kubectl", "-n", "kube-system", "exec", "$(kubectl -n kube-system get pods -l k8s-app='cilium' -o jsonpath='{.items[0].metadata.name}')", "cilium", "config", "PolicyEnforcement=never")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println("Network Policy disabled.")
}

// CreateDockerRegistrySecret ...
func CreateDockerRegistrySecret(name string, dockerServer string, dockerUsername string, dockerPassword string, dockerEmail string) {
	cmd := exec.Command("kubectl", "create", "secret", "docker-registry", name, dockerServer, dockerUsername, dockerPassword, dockerEmail)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// DeleteSecret ...
func DeleteSecret(name string) {
	cmd := exec.Command("kubectl", "delete", "secret", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(dst string, kubectlVersion string) {
	if _, err := os.Stat(gokube.GetBinDir() + "/kubectl.exe"); os.IsNotExist(err) {
		download.DownloadFromUrl("kubectl "+kubectlVersion, URL, kubectlVersion)
		utils.MoveFile(gokube.GetTempDir()+"/kubectl.exe", dst+"/kubectl.exe")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeleteExecutable ...
func DeleteExecutable() {
	utils.RemoveFile(gokube.GetBinDir() + "/kubectl.exe")
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() {
	utils.CleanDir(utils.GetUserHome() + "/.kube")
}
