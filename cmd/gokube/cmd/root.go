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

// TODO add support for wsl2 driver

package cmd

import (
	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/helmimage"
	"github.com/gemalto/gokube/pkg/helmpush"
	"github.com/gemalto/gokube/pkg/helmspray"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/stern"
	"github.com/spf13/cobra"
	"os"
)

const (
	DEFAULT_KUBERNETES_VERSION  = "v1.22.12"
	DEFAULT_KUBECTL_VERSION     = "v1.22.12"
	DEFAULT_MINIKUBE_VERSION    = "v1.26.1"
	DEFAULT_MINIKUBE_URL        = "https://storage.googleapis.com/minikube/releases/%s/minikube-windows-amd64.exe"
	DEFAULT_MINIKUBE_MEMORY     = 8192
	DEFAULT_MINIKUBE_CPUS       = 4
	DEFAULT_MINIKUBE_DNS_DOMAIN = "cluster.local"
	DEFAULT_DOCKER_VERSION      = "20.10.14"
	DEFAULT_HELM_VERSION        = "v3.9.2"
	DEFAULT_HELM_SPRAY_VERSION  = "v4.0.10"
	DEFAULT_HELM_SPRAY_URL      = "https://github.com/ThalesGroup/helm-spray/releases/download/%s/helm-spray-windows-amd64.tar.gz"
	DEFAULT_HELM_IMAGE_VERSION  = "v1.0.7"
	DEFAULT_HELM_IMAGE_URL      = "https://github.com/cvila84/helm-image/releases/download/%s/helm-image-windows-amd64.tar.gz"
	DEFAULT_HELM_PUSH_VERSION   = "0.10.3"
	DEFAULT_HELM_PUSH_URL       = "https://github.com/chartmuseum/helm-push/releases/download/v%s/helm-push_%s_windows_amd64.tar.gz"
	DEFAULT_STERN_VERSION       = "1.21.0"
	DEFAULT_MINIAPPS_REPO       = "https://thalesgroup.github.io/miniapps"
)

var minikubeURL string
var minikubeVersion string
var dockerVersion string
var kubectlVersion string
var helmVersion string
var helmSprayURL string
var helmSprayVersion string
var helmImageURL string
var helmImageVersion string
var helmPushVersion string
var sternVersion string
var askForUpgrade bool
var verbose bool
var quiet bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gokube",
	Short: `gokube is a nice installer to provide an environment for developing day-to-day with kubernetes & helm on your laptop.`,
	Long:  `gokube is a nice installer to provide an environment for developing day-to-day with kubernetes & helm on your laptop.`,
}

func upgradeHelmPlugins() {
	// TODO rely on helm plugin install
	helmspray.DeletePlugin()
	helmspray.InstallPlugin(helmSprayURL, helmSprayVersion)
	helmimage.DeletePlugin()
	helmimage.InstallPlugin(helmImageURL, helmImageVersion)
	helmpush.DeletePlugin()
	helmpush.InstallPlugin(DEFAULT_HELM_PUSH_URL, helmPushVersion)
}

func upgradeDependencies() {
	minikube.DeleteExecutable()
	minikube.DownloadExecutable(gokube.GetBinDir(), minikubeURL, minikubeVersion)
	helm.DeleteExecutable()
	helm.DownloadExecutable(gokube.GetBinDir(), helmVersion)
	docker.DeleteExecutable()
	docker.DownloadExecutable(gokube.GetBinDir(), dockerVersion)
	kubectl.DeleteExecutable()
	kubectl.DownloadExecutable(gokube.GetBinDir(), kubectlVersion)
	stern.DeleteExecutable()
	stern.DownloadExecutable(gokube.GetBinDir(), sternVersion)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Activate verbose logging")
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
