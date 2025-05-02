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
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/helmimage"
	"github.com/gemalto/gokube/pkg/helmpush"
	"github.com/gemalto/gokube/pkg/helmspray"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/stern"
	"github.com/gemalto/gokube/pkg/k9s"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

const (
	DEFAULT_KUBERNETES_VERSION         = "v1.32.0"
	DEFAULT_KUBECTL_VERSION            = "v1.31.0"
	DEFAULT_MINIKUBE_VERSION           = "v1.35.0"
	DEFAULT_MINIKUBE_MEMORY            = 12288
	DEFAULT_MINIKUBE_CPUS              = 6
	DEFAULT_MINIKUBE_SWAP              = 0
	DEFAULT_MINIKUBE_DISK              = "20g"
	DEFAULT_MINIKUBE_DNS_DOMAIN        = "cluster.local"
	DEFAULT_MINIKUBE_CONTAINER_RUNTIME = "docker"
	DEFAULT_DOCKER_VERSION             = "28.0.4"
	DEFAULT_HELM_VERSION               = "v3.17.0"
	DEFAULT_HELM_SPRAY_VERSION         = "v4.0.13"
	DEFAULT_HELM_IMAGE_VERSION         = "v1.1.0"
	DEFAULT_HELM_PUSH_VERSION          = "0.10.4"
	DEFAULT_STERN_VERSION              = "1.32.0"
	DEFAULT_K9S_VERSION                = "0.50.4"
	DEFAULT_MINIAPPS_REPO              = "https://thalesgroup.github.io/miniapps"
	DEFAULT_GOKUBE_CHECK_IP            = "192.168.99.100"
	DEFAULT_GOKUBE_CIDR                = "192.168.99.1/24"
)

var kubernetesVersion string
var containerRuntime string
var kubectlURL string
var kubectlVersion string
var minikubeURL string
var minikubeVersion string
var dockerURL string
var dockerVersion string
var helmURL string
var helmVersion string
var helmSprayURL string
var helmSprayVersion string
var helmImageURL string
var helmImageVersion string
var helmPushURL string
var helmPushVersion string
var sternURL string
var sternVersion string
var k9sURL string
var k9sVersion string
var askForUpgrade bool
var snapshotName string
var verbose bool
var quiet bool
var force bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gokube",
	Short: `gokube is a nice installer to provide an environment for developing day-to-day with kubernetes & helm on your laptop.`,
	Long:  `gokube is a nice installer to provide an environment for developing day-to-day with kubernetes & helm on your laptop.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Activate verbose logging")
}

func checkMinimumRequirements() {
	if semver.New(kubernetesVersion[1:]).Compare(*semver.New("1.20.0")) < 0 {
		fmt.Println("Error: This gokube version is only compatible with kubernetes version >= 1.20.0")
		os.Exit(1)
	}
	if semver.New(minikubeVersion[1:]).Compare(*semver.New("1.25.0")) < 0 {
		fmt.Println("Error: This gokube version is only compatible with minikube version >= 1.25.0")
		os.Exit(1)
	}
	if semver.New(helmVersion[1:]).Compare(*semver.New("3.0.0-0")) < 0 {
		fmt.Println("Error: This gokube version is only compatible with helm version >= 3.0.0-0")
		os.Exit(1)
	}
	if semver.New(helmSprayVersion[1:]).Compare(*semver.New("4.0.0-0")) < 0 {
		fmt.Println("Error: This gokube version is only compatible with helm-spray version >= 4.0.0-0")
		os.Exit(1)
	}
	if semver.New(helmPushVersion).Compare(*semver.New("0.10.0")) <= 0 {
		fmt.Println("Error: This gokube version is only compatible with helm-push version >= 0.10.0")
		os.Exit(1)
	}
}

func loadURLVersionsFromEnv() {
	kubectlURL = utils.GetValueFromEnv("KUBECTL_URL", kubectl.DEFAULT_URL)
	kubectlVersion = utils.GetValueFromEnv("KUBECTL_VERSION", DEFAULT_KUBECTL_VERSION)
	miniappsRepo = utils.GetValueFromEnv("MINIAPPS_URL", DEFAULT_MINIAPPS_REPO)
	minikubeURL = utils.GetValueFromEnv("MINIKUBE_URL", minikube.DEFAULT_URL)
	minikubeVersion = utils.GetValueFromEnv("MINIKUBE_VERSION", DEFAULT_MINIKUBE_VERSION)
	dockerURL = utils.GetValueFromEnv("DOCKER_URL", docker.DEFAULT_URL)
	dockerVersion = utils.GetValueFromEnv("DOCKER_VERSION", DEFAULT_DOCKER_VERSION)
	helmURL = utils.GetValueFromEnv("HELM_URL", helm.DEFAULT_URL)
	helmVersion = utils.GetValueFromEnv("HELM_VERSION", DEFAULT_HELM_VERSION)
	helmSprayURL = utils.GetValueFromEnv("HELM_SPRAY_URL", helmspray.DEFAULT_URL)
	helmSprayVersion = utils.GetValueFromEnv("HELM_SPRAY_VERSION", DEFAULT_HELM_SPRAY_VERSION)
	helmImageURL = utils.GetValueFromEnv("HELM_IMAGE_URL", helmimage.DEFAULT_URL)
	helmImageVersion = utils.GetValueFromEnv("HELM_IMAGE_VERSION", DEFAULT_HELM_IMAGE_VERSION)
	helmPushURL = utils.GetValueFromEnv("HELM_PUSH_URL", helmpush.DEFAULT_URL)
	helmPushVersion = utils.GetValueFromEnv("HELM_PUSH_VERSION", DEFAULT_HELM_PUSH_VERSION)
	sternURL = utils.GetValueFromEnv("STERN_URL", stern.DEFAULT_URL)
	sternVersion = utils.GetValueFromEnv("STERN_VERSION", DEFAULT_STERN_VERSION)
	k9sURL = utils.GetValueFromEnv("K9S_URL", k9s.DEFAULT_URL)
	k9sVersion = utils.GetValueFromEnv("K9S_VERSION", DEFAULT_K9S_VERSION)
}

func upgradeDependencies() error {
	return gokube.UpgradeDependencies(&gokube.Dependencies{
		MinikubeURL:     minikubeURL,
		MinikubeVersion: minikubeVersion,
		HelmURL:         helmURL,
		HelmVersion:     helmVersion,
		DockerURL:       dockerURL,
		DockerVersion:   dockerVersion,
		KubectlURL:      kubectlURL,
		KubectlVersion:  kubectlVersion,
		SternURL:        sternURL,
		SternVersion:    sternVersion,
		K9sURL:          k9sURL,
		K9sVersion:      k9sVersion,
	})
}

func upgradeHelmPlugins() error {
	return gokube.UpgradeHelmPlugins(&gokube.HelmPlugins{
		SprayURL:     helmSprayURL,
		SprayVersion: helmSprayVersion,
		ImageURL:     helmImageURL,
		ImageVersion: helmImageVersion,
		PushURL:      helmPushURL,
		PushVersion:  helmPushVersion,
	})
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
