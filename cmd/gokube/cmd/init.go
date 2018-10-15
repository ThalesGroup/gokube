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
	"os"

	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"

	"github.com/spf13/cobra"
)

var memory int16
var nCPUs int16
var diskSize string
var httpProxy string
var httpsProxy string
var noProxy string
var upgrade bool
var insecureRegistry string
var minikubeFork string
var minikubeVersion string
var helmVersion string
var kubernetesVersion string

// initCmd represents the start command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)",
	Long:  "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)",
	Run:   initRun,
}

func init() {
	initCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", "v0.30.0", "The minikube version (ex: v0.28.0)")
	initCmd.Flags().StringVarP(&helmVersion, "helm-version", "", "v2.11.0", "The helm version (ex: v2.10.0)")
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", "v1.10.8", "The kubernetes version (ex: v1.10.8)")
	initCmd.Flags().StringVarP(&minikubeFork, "minikube-fork", "", "minikube", "The minikube fork which will be used instead of the official one")
	initCmd.Flags().Int16VarP(&memory, "memory", "m", int16(8192), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&nCPUs, "nCPUs", "c", int16(4), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().StringVarP(&diskSize, "disk-size", "d", "20g", "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", "", "HTTP proxy for minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", "", "HTTPS proxy for minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", "", "No proxy for minikube VM")
	initCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false, "Upgrade if Go Kube! is already installed")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", "", "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	RootCmd.AddCommand(initCmd)
}

func initRun(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "usage: gokube init")
		os.Exit(1)
	}

	if upgrade {
		minikube.Delete()
		minikube.Purge()
		helm.Purge()
		kubectl.Purge()
		docker.Purge()
	}

	// Download dependencies
	minikube.Download(gokube.GetBinDir(), minikubeFork, minikubeVersion)
	helm.Download(gokube.GetBinDir(), helmVersion)
	docker.Download(gokube.GetBinDir())
	kubectl.Download(gokube.GetBinDir())

	// Create virtual machine (minikube)
	minikube.Start(memory, nCPUs, diskSize, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion)

	// Disbale notification for updates
	minikube.ConfigSet("WantUpdateNotification", "false")

	// Switch context to minikube for kubectl and helm
	kubectl.ConfigUseContext("minikube")

	// Install helm
	helm.Init()

	// Add Helm repository
	helm.RepoAdd("monocular", "https://helm.github.io/monocular")
	helm.RepoAdd("miniapps", "https://gemalto.github.io/miniapps")
	helm.RepoUpdate()

	// Deploy Monocular
	helm.UpgradeWithConfiguration("nginx", "kube-system", "controller.hostNetwork=true", "stable/nginx-ingress", "0.25.1")
	helm.UpgradeWithConfiguration("gokube", "kube-system", "api.config.repos[0].name=miniapps,api.config.repos[0].url=https://gemalto.github.io/miniapps,api.config.repos[0].source=https://github.com/gemalto/miniapps/tree/master/charts,api.replicaCount=1,api.image.pullPolicy=IfNotPresent,api.config.cacheRefreshInterval=60,ui.replicaCount=1,ui.image.pullPolicy=IfNotPresent,ui.appName=GoKube,prerender.replicaCount=1,prerender.image.pullPolicy=IfNotPresent", "monocular/monocular", "0.6.3")

	// Configure proxy for Monocular
	kubectl.Patch("kube-system", "deployment", "gokube-monocular-api", "{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"monocular\",\"env\":[{\"name\":\"HTTP_PROXY\",\"value\":\""+httpsProxy+"\"}]}]}}}}")
	kubectl.Patch("kube-system", "deployment", "gokube-monocular-api", "{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"monocular\",\"env\":[{\"name\":\"HTTPS_PROXY\",\"value\":\""+httpProxy+"\"}]}]}}}}")
	kubectl.Patch("kube-system", "deployment", "gokube-monocular-api", "{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"monocular\",\"env\":[{\"name\":\"NO_PROXY\",\"value\":\""+noProxy+"\"}]}]}}}}")

	fmt.Println("\nGoKube has been installed.")
	fmt.Println("Now, you need more or less 10 minutes for running pods...")
	fmt.Println("\nTo verify that pods are running, execute:")
	fmt.Println("> kubectl get pods --all-namespaces")
	fmt.Println("")

}
