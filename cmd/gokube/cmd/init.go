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
	"github.com/gemalto/gokube/pkg/utils"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"

	"github.com/spf13/cobra"
)

const (
	NGINX_INGRESS_APP_VERSION  = "0.23.0"
	TPROXY_CHART_VERSION       = "1.0.0"
	DEFAULT_KUBERNETES_VERSION = "v1.10.13"
	DEFAULT_MINIKUBE_VERSION   = "v1.1.0"
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
var upgrade bool
var clean bool
var imageCache bool
var imageCacheAlternateRepo string
var miniappsRepo string
var ingressController bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates the virtual machine (minikube)",
	Long:  "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates the virtual machine (minikube)",
	Run:   initRun,
}

func init() {
	var defaultKubernetesVersion = os.Getenv("KUBERNETES_VERSION")
	if len(defaultKubernetesVersion) == 0 {
		defaultKubernetesVersion = DEFAULT_KUBERNETES_VERSION
	}
	var defaultMinikubeVersion = os.Getenv("MINIKUBE_VERSION")
	if len(defaultMinikubeVersion) == 0 {
		defaultMinikubeVersion = DEFAULT_MINIKUBE_VERSION
	}
	initCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", "https://storage.googleapis.com/minikube/releases/%s/minikube-windows-amd64.exe", "The URL to download minikube")
	initCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", defaultMinikubeVersion, "The minikube version")
	initCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", "18.09.0", "The docker version")
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", defaultKubernetesVersion, "The kubernetes version")
	initCmd.Flags().StringVarP(&kubectlVersion, "kubectl-version", "", "v1.13.6", "The kubectl version")
	initCmd.Flags().StringVarP(&helmVersion, "helm-version", "", "v2.13.1", "The helm version")
	initCmd.Flags().StringVarP(&helmSprayVersion, "helm-spray-version", "", "v3.4.2", "The helm version")
	initCmd.Flags().StringVarP(&sternVersion, "stern-version", "", "1.10.0", "The stern version")
	initCmd.Flags().Int16VarP(&memory, "memory", "", int16(8192), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&cpus, "cpus", "", int16(4), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().StringVarP(&disk, "disk", "", "20g", "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().StringVarP(&checkIP, "check-ip", "", "192.168.99.100", "Checks if minikube VM allocated IP matches the provided one (0.0.0.0 means no check)")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", os.Getenv("INSECURE_REGISTRY"), "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", os.Getenv("HTTP_PROXY"), "HTTP proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", os.Getenv("HTTPS_PROXY"), "HTTPS proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", os.Getenv("NO_PROXY"), "No proxy variable for docker engine in minikube VM")
	initCmd.Flags().BoolVarP(&transparentProxy, "transparent-proxy", "", false, "Manage HTTP proxy connections with transparent proxy, implies --image-cache")
	initCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	initCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Clean gokube (remove docker, minikube, kubectl and helm working directories)")
	initCmd.Flags().BoolVarP(&imageCache, "image-cache", "", true, "Download docker images in cache before pulling them in minikube")
	initCmd.Flags().StringVarP(&imageCacheAlternateRepo, "image-cache-alternate-repo", "", os.Getenv("ALTERNATE_REPO"), "Alternate docker repo used to download images in cache")
	initCmd.Flags().StringVarP(&miniappsRepo, "miniapps-repo", "", "https://gemalto.github.io/miniapps", "Helm repository for miniapps")
	initCmd.Flags().BoolVarP(&ingressController, "ingress-controller", "", false, "Deploy ingress controller")
	RootCmd.AddCommand(initCmd)
}

// Because of https://github.com/google/go-containerregistry/issues/119, we cannot directly cache images from quay.io repository. Temporary fix is to pull the image elsewhere (docker.io) and tag it again with quay.io
func cache(originalImagePath string, alternateImagePath string, imageName string, dockerEnv []utils.EnvVar) {
	var originalImage = ""
	if originalImagePath != "" {
		originalImage = originalImagePath + "/" + imageName
	} else {
		originalImage = imageName
	}
	if imageCacheAlternateRepo != "" && strings.HasPrefix(originalImagePath, "quay.io") {
		var alternateImage = alternateImagePath + "/" + imageName
		minikube.Cache(alternateImage)
		docker.TagImage(alternateImage, originalImage, dockerEnv)
		docker.RemoveImage(alternateImage, dockerEnv)
	} else {
		minikube.Cache(originalImage)
	}
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

	if clean {
		fmt.Println("Deleting gokube working directories...")
		minikube.DeleteWorkingDirectory()
		helm.DeleteWorkingDirectory()
		kubectl.DeleteWorkingDirectory()
		docker.DeleteWorkingDirectory()
		docker.InitWorkingDirectory()
	}
	if upgrade {
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
	// Disable notification for updates
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
		fmt.Println("Caching additional docker images...")
		dockerEnv := minikube.DockerEnv()

		// Put needed images in cache (Helm)
		cache("gcr.io/kubernetes-helm", imageCacheAlternateRepo, "tiller:"+helmVersion, dockerEnv)

		if ingressController {
			// Put needed images in cache (Nginx ingress controller)
			cache("quay.io/kubernetes-ingress-controller", imageCacheAlternateRepo, "nginx-ingress-controller:"+NGINX_INGRESS_APP_VERSION, dockerEnv)
			cache("gcr.io/google_containers", imageCacheAlternateRepo, "defaultbackend:1.4", dockerEnv)
		}

		if transparentProxy && httpProxy != "" && httpsProxy != "" {
			// Put needed images in cache (any-proxy)
			cache("", imageCacheAlternateRepo, "alpine:3.8", dockerEnv)
			cache("", imageCacheAlternateRepo, "cvila84/any-proxy:1.0.1", dockerEnv)
		}
	}

	// Switch context to minikube for kubectl and helm
	kubectl.ConfigUseContext("minikube")

	// Install helm
	fmt.Println("Initializing helm...")
	helm.Init()

	// Add helm repository
	helm.RepoAdd("miniapps", miniappsRepo)
	helm.RepoUpdate()

	if upgrade {
		// Add helm spray plugin
		helmspray.DeletePlugin()
		helmspray.InstallPlugin(helmSprayVersion)
	}

	if ingressController {
		// Deploy ingress controller
		fmt.Println("Deploying ingress controller...")
		minikube.AddonsEnable("ingress")
		time.Sleep(60 * time.Second)
		kubectl.Patch("kube-system", "deploy", "nginx-ingress-controller", "{\"spec\": {\"template\": {\"spec\": {\"hostNetwork\": true}}}}")
		//	helm.UpgradeWithConfiguration("nginx", "kube-system", "controller.hostNetwork=true", "stable/nginx-ingress", NGINX_INGRESS_CHART_VERSION)
	}

	if transparentProxy && httpProxy != "" && httpsProxy != "" {
		// Deploy transparent proxy
		fmt.Println("Installing transparent proxy...")
		helm.UpgradeWithConfiguration("any-proxy", "kube-system", "global.httpProxy="+httpProxy+",global.httpsProxy="+httpsProxy, "miniapps/any-proxy", TPROXY_CHART_VERSION)
	}

	// Patch kubernetes-dashboard to expose it on nodePort 30000
	fmt.Print("Exposing kubernetes dashboard...")
	for n := 1; n < 12; n++ {
		var dashboardService = kubectl.GetObject("kube-system", "svc", "kubernetes-dashboard")
		if len(dashboardService) > 0 {
			fmt.Println()
			kubectl.Patch("kube-system", "svc", "kubernetes-dashboard", "{\"spec\":{\"type\":\"NodePort\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":9090,\"nodePort\":30000}]}}")
			break
		} else {
			fmt.Print(".")
			time.Sleep(10 * time.Second)
		}
	}

	fmt.Println("\ngokube has been installed.")
	if !imageCache {
		fmt.Println("Now, you need more or less 10 minutes for running pods...")
	}
	fmt.Println("\nTo verify that pods are running, execute:")
	fmt.Println("> kubectl get pods --all-namespaces")
	fmt.Println("")
}
