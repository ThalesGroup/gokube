package helmpush

import (
	"github.com/coreos/go-semver/semver"
	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
	"os"
)

// InstallPlugin ...
func InstallPlugin(helmPushURI string, helmPushVersion string) {
	var binaryName string
	if semver.New(helmPushVersion).Compare(*semver.New("0.9.0")) <= 0 {
		binaryName = "/bin/helmpush.exe"
	} else {
		binaryName = "/bin/helm-cm-push.exe"
	}

	var helm3PluginHome = utils.GetAppDataHome() + "/helm/plugins/helm-push"
	if _, err := os.Stat(helm3PluginHome + binaryName); os.IsNotExist(err) {
		download.FromUrl("helm-push "+helmPushVersion, helmPushURI, helmPushVersion)
		utils.CreateDirs(helm3PluginHome + "/bin")
		utils.MoveFile(gokube.GetTempDir()+binaryName, helm3PluginHome+binaryName)
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
