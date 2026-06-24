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
	"path/filepath"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	DEFAULT_URL           = "https://dl.k8s.io/%s/bin/windows/amd64/kubectl.exe"
	LOCAL_EXECUTABLE_NAME = "kubectl.exe"
)

func bin() string {
	return utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
}

// Get ...
func Get(namespace string, resourceType string, resourceName string, jsonPath string) (string, error) {
	var args = []string{"--namespace", namespace, "get", resourceType, resourceName}
	if len(jsonPath) > 0 {
		args = append(args, "-o", "jsonpath="+jsonPath)
	}
	output, err := exec.Command(bin(), args...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// ConfigUseContext ...
func ConfigUseContext(context string) error {
	cmd := exec.Command(bin(), "config", "use-context", context)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Patch ...
func Patch(namespace string, resourceType string, resourceName string, patch string) error {
	cmd := exec.Command(bin(), "--namespace", namespace, "patch", resourceType, resourceName, "-p", patch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Version ...
func Version() error {
	fmt.Println("kubectl version: ")
	cmd := exec.Command(bin(), "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DownloadExecutable ...
func DownloadExecutable(kubectlURL string, kubectlVersion string, bar *pb.ProgressBar) error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if download.IsCurrentVersion(localFile, kubectlVersion) {
		bar.SetTemplateString(`{{ green "kubectl" }} ` + kubectlVersion + ` already up to date`)
		bar.SetTotal(1)
		bar.SetCurrent(1)
		bar.Finish()
		return nil
	}
	_ = os.RemoveAll(localFile)
	_ = os.RemoveAll(download.VersionFile(localFile))
	fileMap := &download.FileMap{Src: LOCAL_EXECUTABLE_NAME, Dst: LOCAL_EXECUTABLE_NAME}
	if _, err := download.FromUrl(kubectlURL, kubectlVersion, "kubectl", []*download.FileMap{fileMap}, filepath.Dir(localFile), bar); err != nil {
		return err
	}
	return download.WriteVersion(localFile, kubectlVersion)
}

// DeleteExecutable ...
func DeleteExecutable() error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	_ = os.RemoveAll(download.VersionFile(localFile))
	return os.RemoveAll(localFile)
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() error {
	return utils.CleanDir(utils.GetUserHome() + "/.kube")
}
