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

package helmspray

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	URL = "https://github.com/cvila84/helm-spray/releases/download/%s/helm-spray-windows-amd64.tar.gz"
)

//Version ...
func Version() {
	fmt.Print("helm-spray version: ")
	cmd := exec.Command("helm", "plugin", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// InstallPlugin ...
func InstallPlugin(helmSprayVersion string) {
	var pluginHome = utils.GetUserHome() + "/.helm/plugins/helm-spray"
	utils.CreateDir(pluginHome)
	if _, err := os.Stat(pluginHome + "/helm-spray.exe"); os.IsNotExist(err) {
		download.DownloadFromUrl("helm-spray "+helmSprayVersion, URL, helmSprayVersion)
		utils.MoveFile(gokube.GetTempDir()+"/helm-spray.exe", pluginHome+"/helm-spray.exe")
		utils.MoveFile(gokube.GetTempDir()+"/plugin.yaml", pluginHome+"/plugin.yaml")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeletePlugin ...
func DeletePlugin() {
	var pluginHome = utils.GetUserHome() + "/.helm/plugins/helm-spray"
	_, err := os.Stat(pluginHome)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(err)
	}
	utils.RemoveFiles(pluginHome + "/*")
	utils.RemoveDir(pluginHome)
}
