package helmpush

import (
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
)

// InstallPlugin ...
func InstallPlugin(helmPushURI string, helmPushVersion string) {
	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-push"
	if _, err := os.Stat(helm3PluginHome + "/bin/helmpush.exe"); os.IsNotExist(err) {
		download.FromUrl("helm-push "+helmPushVersion, helmPushURI, helmPushVersion)
		utils.CreateDirs(helm3PluginHome + "/bin")
		utils.MoveFile(gokube.GetTempDir()+"/bin/helmpush.exe", helm3PluginHome+"/bin/helmpush.exe")
		utils.MoveFile(gokube.GetTempDir()+"/plugin.yaml", helm3PluginHome+"/plugin.yaml")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeletePlugin ...
func DeletePlugin() {
	var helm2PluginHome = utils.GetUserHome() + "/.helm/plugins/helm-push"
	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-push"
	utils.RemoveDir(helm2PluginHome)
	utils.RemoveDir(helm3PluginHome)
}
