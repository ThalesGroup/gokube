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
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:          "start",
	Short:        "Starts gokube. This command starts minikube",
	Long:         "Starts gokube. This command starts minikube",
	RunE:         startRun,
	SilenceUsage: true,
}

func init() {
	var defaultKubectlVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBECTL_VERSION)
	var defaultMinikubeUrl = getValueFromEnv("MINIKUBE_URL", DEFAULT_MINIKUBE_URL)
	var defaultMinikubeVersion = getValueFromEnv("MINIKUBE_VERSION", DEFAULT_MINIKUBE_VERSION)
	var defaultDockerVersion = getValueFromEnv("DOCKER_VERSION", DEFAULT_DOCKER_VERSION)
	var defaultHelmVersion = getValueFromEnv("HELM_VERSION", DEFAULT_HELM_VERSION)
	var defaultHelmSprayUrl = getValueFromEnv("HELM_SPRAY_URL", DEFAULT_HELM_SPRAY_URL)
	var defaultHelmSprayVersion = getValueFromEnv("HELM_SPRAY_VERSION", DEFAULT_HELM_SPRAY_VERSION)
	var defaultHelmImageUrl = getValueFromEnv("HELM_IMAGE_URL", DEFAULT_HELM_IMAGE_URL)
	var defaultHelmImageVersion = getValueFromEnv("HELM_IMAGE_VERSION", DEFAULT_HELM_IMAGE_VERSION)
	var defaultHelmPushVersion = getValueFromEnv("HELM_PUSH_VERSION", DEFAULT_HELM_PUSH_VERSION)
	startCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", defaultMinikubeUrl, "The URL to download minikube")
	startCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", defaultMinikubeVersion, "The minikube version")
	startCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", defaultDockerVersion, "The docker version")
	startCmd.Flags().StringVarP(&kubectlVersion, "kubectl-version", "", defaultKubectlVersion, "The kubectl version")
	startCmd.Flags().StringVarP(&helmVersion, "helm-version", "", defaultHelmVersion, "The helm version")
	startCmd.Flags().StringVarP(&helmSprayURL, "helm-spray-url", "", defaultHelmSprayUrl, "The URL to download helm spray plugin")
	startCmd.Flags().StringVarP(&helmSprayVersion, "helm-spray-version", "", defaultHelmSprayVersion, "The helm spray plugin version")
	startCmd.Flags().StringVarP(&helmImageURL, "helm-image-url", "", defaultHelmImageUrl, "The URL to download helm image plugin")
	startCmd.Flags().StringVarP(&helmImageVersion, "helm-image-version", "", defaultHelmImageVersion, "The helm image plugin version")
	startCmd.Flags().StringVarP(&helmPushVersion, "helm-push-version", "", defaultHelmPushVersion, "The helm push plugin version")
	startCmd.Flags().StringVarP(&sternVersion, "stern-version", "", DEFAULT_STERN_VERSION, "The stern version")
	startCmd.Flags().BoolVarP(&askForUpgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	rootCmd.AddCommand(startCmd)
}

func start() error {
	gokube.ReadConfig(verbose)
	kubernetesVersionForStart := viper.GetString("kubernetes-version")
	if len(kubernetesVersionForStart) == 0 {
		kubernetesVersionForStart = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBERNETES_VERSION)
	}
	fmt.Printf("Starting minikube VM with kubernetes %s...\n", kubernetesVersion)
	minikube.Restart(kubernetesVersion)
	return nil
}

func startRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}
	if askForUpgrade {
		fmt.Println("Upgrading gokube dependencies...")
		upgradeDependencies()
		fmt.Println("Upgrading helm plugins...")
		upgradeHelmPlugins()
	}
	return start()
}
