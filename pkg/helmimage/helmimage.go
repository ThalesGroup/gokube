package helmimage

import (
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
)

// InstallPlugin ...
func InstallPlugin(helmImageURI string, helmImageVersion string) {
	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-image"
	if _, err := os.Stat(helm3PluginHome + "/bin/helm-image.exe"); os.IsNotExist(err) {
		download.FromUrl("helm-image "+helmImageVersion, helmImageURI, helmImageVersion)
		utils.CreateDirs(helm3PluginHome + "/bin")
		utils.MoveFile(gokube.GetTempDir()+"/bin/helm-image.exe", helm3PluginHome+"/bin/helm-image.exe")
		utils.MoveFile(gokube.GetTempDir()+"/plugin.yaml", helm3PluginHome+"/plugin.yaml")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeletePlugin ...
func DeletePlugin() {
	var helm2PluginHome = utils.GetUserHome() + "/.helm/plugins/helm-image"
	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-image"
	utils.RemoveDir(helm2PluginHome)
	utils.RemoveDir(helm3PluginHome)
}
