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

package helmimage

import (
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
)

var (
	DEFAULT_URL           = defaultURL()
	LOCAL_EXECUTABLE_NAME = utils.ExecutableName("helm-image")
)

func defaultURL() string {
	if runtime.GOOS == "linux" {
		return "https://github.com/ThalesGroup/helm-image/releases/download/%s/helm-image-linux-amd64.tar.gz"
	}
	return "https://github.com/ThalesGroup/helm-image/releases/download/%s/helm-image-windows-amd64.tar.gz"
}

// InstallPlugin ...
func InstallPlugin(helmImageURI string, helmImageVersion string) error {
	pluginDir := utils.GetAppDataHome() + string(os.PathSeparator) +
		"helm" + string(os.PathSeparator) +
		"plugins" + string(os.PathSeparator) +
		"helm-image"
	localFile := pluginDir + string(os.PathSeparator) +
		"bin" + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		fileMap1 := &download.FileMap{Src: "bin" + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME, Dst: "bin" + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME}
		containerdName := utils.ExecutableName("containerd")
		fileMap2 := &download.FileMap{Src: "bin" + string(os.PathSeparator) + containerdName, Dst: "bin" + string(os.PathSeparator) + containerdName}
		fileMap3 := &download.FileMap{Src: "plugin.yaml", Dst: "plugin.yaml"}
		_, err = download.FromUrl(helmImageURI, helmImageVersion, "helm-image", []*download.FileMap{fileMap1, fileMap2, fileMap3}, pluginDir)
		if err != nil {
			return err
		}
		if err := utils.MakeExecutable(localFile); err != nil {
			return err
		}
		return utils.MakeExecutable(filepath.Join(filepath.Dir(localFile), containerdName))
	}
	return nil
}

// DeletePlugin ...
func DeletePlugin() error {
	localFile := utils.GetAppDataHome() + string(os.PathSeparator) +
		"helm" + string(os.PathSeparator) +
		"plugins" + string(os.PathSeparator) +
		"helm-image" + string(os.PathSeparator)
	return os.RemoveAll(localFile)
}
