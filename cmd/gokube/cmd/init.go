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
	"github.com/gemalto/gokube/pkg/helmspray"
	"github.com/gemalto/gokube/pkg/stern"
	"log"
	"os"
	"strings"

	"github.com/gemalto/gokube/pkg/utils"

	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"

	"github.com/spf13/cobra"
)

const (
	MONOCULAR_CHART_VERSION     = "1.2.8"
	MONOCULAR_APP_VERSION       = "1.2.0"
	NGINX_INGRESS_CHART_VERSION = "1.1.4"
	NGINX_INGRESS_APP_VERSION   = "0.21.0"
	TPROXY_CHART_VERSION        = "1.0.0"
)

var minikubeURL string
var minikubeVersion string
var dockerVersion string
var kubernetesVersion string
var kubectlVersion string
var helmVersion string
var helmSprayVersion string
var sternVersion string
var memory int16
var cpus int16
var disk string
var checkIP string
var insecureRegistry string
var httpProxy string
var httpsProxy string
var noProxy string
var transparentProxy bool
var installMonocular bool
var upgrade bool
var clean bool
var imageCache bool
var imageCacheAlternateRepo string
var miniappsRepo string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)",
	Long:  "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)",
	Run:   initRun,
}

func init() {
	initCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", "https://storage.googleapis.com/minikube/releases/%s/minikube-windows-amd64.exe", "The URL to download minikube")
	initCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", "v0.33.1", "The minikube version")
	initCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", "18.06.1-ce", "The docker version")
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", "v1.10.12", "The kubernetes version")
	initCmd.Flags().StringVarP(&kubectlVersion, "kubectl-version", "", "v1.13.2", "The kubectl version")
	initCmd.Flags().StringVarP(&helmVersion, "helm-version", "", "v2.12.2", "The helm version")
	initCmd.Flags().StringVarP(&helmSprayVersion, "helm-spray-version", "", "v3.2.0", "The helm version")
	initCmd.Flags().StringVarP(&sternVersion, "stern-version", "", "1.10.0", "The stern version")
	initCmd.Flags().Int16VarP(&memory, "memory", "", int16(8192), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&cpus, "cpus", "", int16(4), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().StringVarP(&disk, "disk", "", "20g", "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().StringVarP(&checkIP, "check-ip", "", "192.168.99.100", "Checks if minikube VM allocated IP matches the provided one (0.0.0.0 means no check)")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", "", "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", os.Getenv("HTTP_PROXY"), "HTTP proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", os.Getenv("HTTPS_PROXY"), "HTTPS proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", os.Getenv("NO_PROXY"), "No proxy variable for docker engine in minikube VM")
	initCmd.Flags().BoolVarP(&transparentProxy, "transparent-proxy", "", false, "Manage HTTP proxy connections with transparent proxy, implies --image-cache")
	initCmd.Flags().BoolVarP(&installMonocular, "install-monocular", "", true, "Install monocular")
	initCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	initCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Clean gokube (remove docker, minikube, kubectl and helm working directories)")
	initCmd.Flags().BoolVarP(&imageCache, "image-cache", "", false, "Download docker images in cache before pulling them in minikube")
	initCmd.Flags().StringVarP(&imageCacheAlternateRepo, "image-cache-alternate-repo", "", "", "Alternate docker repo used to download images in cache")
	initCmd.Flags().StringVarP(&miniappsRepo, "miniapps-repo", "", "https://gemalto.github.io/miniapps", "Helm repository for miniapps")
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

	// TODO add manifest to ask for admin rights
	fmt.Println("Deleting previous minikube VM...")
	minikube.Delete()
	//  Does not work well with VB6 and not yet tested with VB5
	//	fmt.Println("Deleting host-only network used by minikube...")
	//	virtualbox.PurgeHostOnlyNetwork()

	if upgrade {
		if clean {
			fmt.Println("Deleting gokube working directories...")
			minikube.DeleteWorkingDirectory()
			helm.DeleteWorkingDirectory()
			kubectl.DeleteWorkingDirectory()
			docker.DeleteWorkingDirectory()
			docker.InitWorkingDirectory()
		}
		fmt.Println("Downloading gokube dependencies...")
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

	if transparentProxy {
		imageCache = true
	}

	// Create virtual machine (minikube)
	minikube.Start(memory, cpus, disk, transparentProxy, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion, imageCache)
	// Disbale notification for updates
	minikube.ConfigSet("WantUpdateNotification", "false")
	// Enable dashboard
	minikube.ConfigSet("dashboard", "true")
	// Checks minikube IP
	var minikubeIP = minikube.Ip()
	if strings.Compare("0.0.0.0", checkIP) != 0 && strings.Compare(checkIP, minikubeIP) != 0 {
		log.Fatalf("Minikube IP (%s) does not match expected IP (%s)", minikubeIP, checkIP)
		os.Exit(1)
	}

	if imageCache {
		// TODO use version const for docker images to match helm versions
		fmt.Println("Putting docker images in cache...")
		dockerEnv := minikube.DockerEnv()

		// Put needed images in cache (Helm)
		minikube.Cache("gcr.io/kubernetes-helm/tiller:" + helmVersion)

		// Put needed images in cache (Nginx ingress controller)
		cacheAndTag(imageCacheAlternateRepo, "nginx-ingress-controller:"+NGINX_INGRESS_APP_VERSION, "quay.io/kubernetes-ingress-controller", dockerEnv)
		minikube.Cache("k8s.gcr.io/defaultbackend:1.4")

		if installMonocular {
			// Put needed images in cache (Monocular)
			cacheAndTag(imageCacheAlternateRepo, "chart-repo:v"+MONOCULAR_APP_VERSION, "quay.io/helmpack", dockerEnv)
			cacheAndTag(imageCacheAlternateRepo, "chartsvc:v"+MONOCULAR_APP_VERSION, "quay.io/helmpack", dockerEnv)
			cacheAndTag(imageCacheAlternateRepo, "monocular-ui:v"+MONOCULAR_APP_VERSION, "quay.io/helmpack", dockerEnv)
			minikube.Cache("docker.io/bitnami/mongodb:4.0.3")
			minikube.Cache("migmartri/prerender:latest")
		}

		if transparentProxy && httpProxy != "" && httpsProxy != "" {
			// Put needed images in cache (any-proxy)
			minikube.Cache("alpine:3.8")
			minikube.Cache(imageCacheAlternateRepo + "/any-proxy:1.0.1")
		}
	}

	// Switch context to minikube for kubectl and helm
	kubectl.ConfigUseContext("minikube")

	// Install helm
	fmt.Println("Initializing helm...")
	helm.Init()

	// Add Helm repository
	if installMonocular {
		helm.RepoAdd("monocular", "https://helm.github.io/monocular")
	}
	helm.RepoAdd("miniapps", miniappsRepo)
	helm.RepoUpdate()

	if upgrade {
		helmspray.DeletePlugin()
		helmspray.InstallPlugin(helmSprayVersion)
	}

	//	minikube.AddonsEnable("ingress")
	helm.UpgradeWithConfiguration("nginx", "kube-system", "controller.hostNetwork=true", "stable/nginx-ingress", NGINX_INGRESS_CHART_VERSION)

	if installMonocular {
		fmt.Println("Installing monocular...")
		var goKubeConfiguration = "sync.repos[0].name=miniapps,sync.repos[0].url=" + miniappsRepo + ",chartsvc.replicas=1,ui.replicaCount=1,ui.image.pullPolicy=IfNotPresent,ui.appName=gokube,prerender.image.pullPolicy=IfNotPresent,ingress.hosts[0]="
		if !transparentProxy && httpProxy != "" && httpsProxy != "" {
			goKubeConfiguration = goKubeConfiguration + ",sync.httpProxy=" + httpProxy + ",sync.httpsProxy=" + httpsProxy
		}
		helm.UpgradeWithConfiguration("gokube", "kube-system", goKubeConfiguration, "monocular/monocular", MONOCULAR_CHART_VERSION)
	}

	// Deploy transparent proxy (if requested)
	if transparentProxy && httpProxy != "" && httpsProxy != "" {
		fmt.Println("Installing transparent proxy...")
		helm.UpgradeWithConfiguration("any-proxy", "kube-system", "global.httpProxy="+httpProxy+",global.httpsProxy="+httpsProxy, "miniapps/any-proxy", TPROXY_CHART_VERSION)
	}

	// Patch kubernetes-dashboard to expose it on nodePort 30000
	kubectl.Patch("kube-system", "svc", "kubernetes-dashboard", "{\"spec\":{\"type\":\"NodePort\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":9090,\"nodePort\":30000}]}}")

	fmt.Println("\ngokube has been installed.")
	if !imageCache {
		fmt.Println("Now, you need more or less 10 minutes for running pods...")
	}
	fmt.Println("\nTo verify that pods are running, execute:")
	fmt.Println("> kubectl get pods --all-namespaces")
	fmt.Println("")
}
