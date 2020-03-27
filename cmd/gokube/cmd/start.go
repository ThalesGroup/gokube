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

package cmd

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts minikube. This command starts minikube",
	Long:  "Starts minikube. This command starts minikube",
	Run:   startRun,
}

func init() {
	gokube.ReadConfig()
	var defaultKubernetesVersion = viper.GetString("kubernetes-version")
	if len(defaultKubernetesVersion) == 0 {
		defaultKubernetesVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBERNETES_VERSION)
	}
	var defaultKubectlVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBECTL_VERSION)
	var defaultMinikubeUrl = getValueFromEnv("MINIKUBE_URL", DEFAULT_MINIKUBE_URL)
	var defaultMinikubeVersion = getValueFromEnv("MINIKUBE_VERSION", DEFAULT_MINIKUBE_VERSION)
	var defaultDockerVersion = getValueFromEnv("DOCKER_VERSION", DEFAULT_DOCKER_VERSION)
	var defaultHelmVersion = getValueFromEnv("HELM_VERSION", DEFAULT_HELM_VERSION)
	var defaultHelmSprayVersion = getValueFromEnv("HELM_SPRAY_VERSION", DEFAULT_HELM_SPRAY_VERSION)
	startCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", defaultMinikubeUrl, "The URL to download minikube")
	startCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", defaultMinikubeVersion, "The minikube version")
	startCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", defaultDockerVersion, "The docker version")
	startCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", defaultKubernetesVersion, "The kubernetes version")
	startCmd.Flags().StringVarP(&kubectlVersion, "kubectl-version", "", defaultKubectlVersion, "The kubectl version")
	startCmd.Flags().StringVarP(&helmVersion, "helm-version", "", defaultHelmVersion, "The helm version")
	startCmd.Flags().StringVarP(&helmSprayVersion, "helm-spray-version", "", defaultHelmSprayVersion, "The helm spray plugin version")
	startCmd.Flags().StringVarP(&sternVersion, "stern-version", "", DEFAULT_STERN_VERSION, "The stern version")
	startCmd.Flags().BoolVarP(&askForUpgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	RootCmd.AddCommand(startCmd)
}

func startRun(cmd *cobra.Command, args []string) {

	helm3 = false
	if semver.New(helmVersion[1:]).Compare(*semver.New("3.0.0")) >= 0 {
		helm3 = true
	}

	if askForUpgrade {
		fmt.Println("Downloading gokube dependencies...")
		upgrade()
		fmt.Println("Installing helm plugins...")
		installHelmPlugins()
	}

	fmt.Printf("Starting minikube VM with kubernetes %s...\n", kubernetesVersion)
	minikube.Restart(kubernetesVersion)
}
