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
	pb "github.com/cheggaaa/pb/v3"
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
)

const (
	DEFAULT_URL           = "https://github.com/ThalesGroup/helm-image/releases/download/%s/helm-image-windows-amd64.tar.gz"
	LOCAL_EXECUTABLE_NAME = "helm-image.exe"
)

// InstallPlugin ...
func InstallPlugin(helmImageURI string, helmImageVersion string, bar *pb.ProgressBar) error {
	sep := string(os.PathSeparator)
	pluginDir := utils.GetAppDataHome() + sep + "helm" + sep + "plugins" + sep + "helm-image"
	installedBinary := pluginDir + sep + "bin" + sep + LOCAL_EXECUTABLE_NAME
	if download.IsCurrentVersion(installedBinary, helmImageVersion) {
		bar.SetTemplateString(`{{ green "helm-image" }} ` + helmImageVersion + ` already up to date`)
		bar.SetTotal(1)
		bar.SetCurrent(1)
		bar.Finish()
		return nil
	}
	_ = os.RemoveAll(pluginDir)
	_ = os.RemoveAll(download.VersionFile(installedBinary))
	fileMap1 := &download.FileMap{Src: "bin" + sep + LOCAL_EXECUTABLE_NAME, Dst: "bin" + sep + LOCAL_EXECUTABLE_NAME}
	fileMap2 := &download.FileMap{Src: "bin" + sep + "containerd.exe", Dst: "bin" + sep + "containerd.exe"}
	fileMap3 := &download.FileMap{Src: "plugin.yaml", Dst: "plugin.yaml"}
	if _, err := download.FromUrl(helmImageURI, helmImageVersion, "helm-image", []*download.FileMap{fileMap1, fileMap2, fileMap3}, pluginDir, bar); err != nil {
		return err
	}
	return download.WriteVersion(installedBinary, helmImageVersion)
}

// DeletePlugin ...
func DeletePlugin() error {
	localFile := utils.GetAppDataHome() + string(os.PathSeparator) +
		"helm" + string(os.PathSeparator) +
		"plugins" + string(os.PathSeparator) +
		"helm-image" + string(os.PathSeparator)
	return os.RemoveAll(localFile)
}
