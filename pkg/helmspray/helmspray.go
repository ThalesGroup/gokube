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

//Version ...
func Version() {
	fmt.Print("helm-spray version:\n")
	cmd := exec.Command("helm", "plugin", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// InstallPlugin ...
func InstallPlugin(helmSprayURI string, helmSprayVersion string) {
	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-spray"
	if _, err := os.Stat(helm3PluginHome + "/bin/helm-spray.exe"); os.IsNotExist(err) {
		download.FromUrl("helm-spray "+helmSprayVersion, helmSprayURI, helmSprayVersion)
		utils.CreateDirs(helm3PluginHome + "/bin")
		utils.MoveFile(gokube.GetTempDir()+"/bin/helm-spray.exe", helm3PluginHome+"/bin/helm-spray.exe")
		utils.MoveFile(gokube.GetTempDir()+"/plugin.yaml", helm3PluginHome+"/plugin.yaml")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeletePlugin ...
func DeletePlugin() {
	var helm2PluginHome = utils.GetUserHome() + "/.helm/plugins/helm-spray"
	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-spray"
	utils.RemoveDir(helm2PluginHome)
	utils.RemoveDir(helm3PluginHome)
}
