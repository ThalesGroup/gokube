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

	"github.com/gemalto/gokube/pkg/utils"

	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"

	"github.com/spf13/cobra"
)

var memory int16
var cpus int16
var diskSize string
var tproxy bool
var httpProxy string
var httpsProxy string
var noProxy string
var upgrade bool
var insecureRegistry string
var minikubeURL string
var minikubeVersion string
var helmVersion string
var kubernetesVersion string
var cache bool
var alternateCacheImagePath string
var miniappsHelmRepository string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)",
	Long:  "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)",
	Run:   initRun,
}

func init() {
	initCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", "v0.32.0", "The minikube version (ex: v0.32.0)")
	initCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", "https://storage.googleapis.com/minikube/releases/%s/minikube-windows-amd64.exe", "The URL to download minikube")
	initCmd.Flags().StringVarP(&helmVersion, "helm-version", "", "v2.12.1", "The helm version (ex: v2.12.1)")
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", "v1.10.12", "The kubernetes version (ex: v1.10.12)")
	initCmd.Flags().Int16VarP(&memory, "memory", "m", int16(8192), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&cpus, "cpus", "c", int16(4), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().StringVarP(&diskSize, "disk-size", "d", "20g", "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().BoolVarP(&tproxy, "transparent-proxy", "", false, "Manage HTTP proxy connections with transparent proxy, implies --cache")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", "", "HTTP proxy for minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", "", "HTTPS proxy for minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", "", "No proxy for minikube VM")
	initCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false, "Upgrade if Go Kube! is already installed")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", "", "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	initCmd.Flags().BoolVarP(&cache, "cache", "", false, "Download images in cache before pulling them in minikube")
	initCmd.Flags().StringVarP(&alternateCacheImagePath, "alternate-cache-image-path", "", "", "Alternate docker image path used to download images in cache")
	initCmd.Flags().StringVarP(&miniappsHelmRepository, "miniapps-helm-repository", "", "https://gemalto.github.io/miniapps", "Helm repository for miniapps")
	RootCmd.AddCommand(initCmd)
}

// Because of https://github.com/google/go-containerregistry/issues/119, we cannot directly cache images from quay.io repository. Temporary fix is to pull the image elsewhere (docker.io) and tag it again with quay.io
func cacheAndTag(imagePath string, imageName string, originalPath string, dockerEnv []utils.EnvVar) {
	var image = imagePath + "/" + imageName
	var originalImage = originalPath + "/" + imageName
	minikube.Cache(image)
	docker.TagImage(image, originalImage, dockerEnv)
	docker.RemoveImage(image, dockerEnv)
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
		docker.Init()
	}

	// Download dependencies...
	minikube.Download(gokube.GetBinDir(), minikubeURL, minikubeVersion)
	helm.Download(gokube.GetBinDir(), helmVersion)
	docker.Download(gokube.GetBinDir())
	kubectl.Download(gokube.GetBinDir())

	if tproxy {
		cache = true
	}

	// Create virtual machine (minikube)
	minikube.Start(memory, cpus, diskSize, tproxy, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion, cache)

	// Disbale notification for updates
	minikube.ConfigSet("WantUpdateNotification", "false")

	// Displays minikube IP
	minikube.Ip()

	if cache {
		dockerEnv := minikube.DockerEnv()

		// Put needed images in cache (Helm)
		minikube.Cache("gcr.io/kubernetes-helm/tiller:v2.11.0")

		// Put needed images in cache (Nginx ingress controller)
		cacheAndTag(alternateCacheImagePath, "nginx-ingress-controller:0.20.0", "quay.io/kubernetes-ingress-controller", dockerEnv)
		minikube.Cache("k8s.gcr.io/defaultbackend:1.4")

		// Put needed images in cache (Monocular)
		cacheAndTag(alternateCacheImagePath, "chart-repo:v1.0.0", "quay.io/helmpack", dockerEnv)
		cacheAndTag(alternateCacheImagePath, "chartsvc:v1.0.0", "quay.io/helmpack", dockerEnv)
		cacheAndTag(alternateCacheImagePath, "monocular-ui:v1.0.0", "quay.io/helmpack", dockerEnv)
		minikube.Cache("docker.io/bitnami/mongodb:4.0.3")
		minikube.Cache("migmartri/prerender:latest")

		if tproxy && httpProxy != "" && httpsProxy != "" {
			// Put needed images in cache (any-proxy)
			minikube.Cache("alpine:3.8")
			minikube.Cache(alternateCacheImagePath + "/any-proxy:1.0.1")
		}
	}

	// Switch context to minikube for kubectl and helm
	kubectl.ConfigUseContext("minikube")

	// Install helm
	helm.Init()

	// Add Helm repository
	helm.RepoAdd("monocular", "https://helm.github.io/monocular")
	helm.RepoAdd("miniapps", miniappsHelmRepository)
	helm.RepoUpdate()

	// Deploy Monocular
	helm.UpgradeWithConfiguration("nginx", "kube-system", "controller.hostNetwork=true", "stable/nginx-ingress", "0.29.2")
	var goKubeConfiguration = "sync.repos[0].name=miniapps,sync.repos[0].url=" + miniappsHelmRepository + ",chartsvc.replicas=1,ui.replicaCount=1,ui.image.pullPolicy=IfNotPresent,ui.appName=gokube,prerender.image.pullPolicy=IfNotPresent"
	if !tproxy && httpProxy != "" && httpsProxy != "" {
		goKubeConfiguration = goKubeConfiguration + ",sync.httpProxy=" + httpProxy + ",sync.httpsProxy=" + httpsProxy
	}
	helm.UpgradeWithConfiguration("gokube", "kube-system", goKubeConfiguration, "monocular/monocular", "1.2.0")

	// Deploy transparent proxy (if requested)
	if tproxy && httpProxy != "" && httpsProxy != "" {
		helm.UpgradeWithConfiguration("any-proxy", "kube-system", "global.httpProxy="+httpProxy+",global.httpsProxy="+httpsProxy, "miniapps/any-proxy", "1.0.0")
	}

	// Patch kubernetes-dashboard to expose it on nodePort 30000
	kubectl.Patch("kube-system", "svc", "kubernetes-dashboard", "{\"spec\":{\"type\":\"NodePort\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":9090,\"nodePort\":30000}]}}")

	fmt.Println("\ngokube has been installed.")
	fmt.Println("Now, you need more or less 10 minutes for running pods...")
	fmt.Println("\nTo verify that pods are running, execute:")
	fmt.Println("> kubectl get pods --all-namespaces")
	fmt.Println("")
}
