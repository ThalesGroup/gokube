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
	"github.com/gemalto/gokube/pkg/helmspray"
	"github.com/gemalto/gokube/pkg/stern"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
)

const (
	DEFAULT_KUBERNETES_VERSION = "v1.18.0"
	DEFAULT_KUBECTL_VERSION    = "v1.18.0"
	DEFAULT_MINIKUBE_VERSION   = "v1.9.1"
	DEFAULT_MINIKUBE_URL       = "https://storage.googleapis.com/minikube/releases/%s/minikube-windows-amd64.exe"
	DEFAULT_DOCKER_VERSION     = "19.03.8"
	DEFAULT_HELM_VERSION       = "v2.16.3"
	DEFAULT_HELM_SPRAY_VERSION = "v3.4.5"
	DEFAULT_STERN_VERSION      = "1.11.0"
)

var gokubeVersion string
var minikubeURL string
var minikubeVersion string
var dockerVersion string
var kubernetesVersion string
var kubectlVersion string
var helmVersion string
var helm3 bool
var helmSprayVersion string
var sternVersion string
var memory int16
var cpus int16
var disk string
var checkIP string
var ipCheckNeeded bool
var insecureRegistry string
var httpProxy string
var httpsProxy string
var noProxy string
var askForUpgrade bool
var askForClean bool
var miniappsRepo string
var dnsProxy bool
var hostDNSResolver bool
var debug bool
var forceInit bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates the virtual machine (minikube)",
	Long:  "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates the virtual machine (minikube)",
	Run:   initRun,
}

func init() {
	gokube.ReadConfig()
	gokubeVersion = viper.GetString("gokube-version")
	if len(gokubeVersion) == 0 {
		gokubeVersion = "0.0.0"
	}
	var defaultKubernetesVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBERNETES_VERSION)
	var defaultKubectlVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBECTL_VERSION)
	var defaultMinikubeUrl = getValueFromEnv("MINIKUBE_URL", DEFAULT_MINIKUBE_URL)
	var defaultMinikubeVersion = getValueFromEnv("MINIKUBE_VERSION", DEFAULT_MINIKUBE_VERSION)
	var defaultDockerVersion = getValueFromEnv("DOCKER_VERSION", DEFAULT_DOCKER_VERSION)
	var defaultHelmVersion = getValueFromEnv("HELM_VERSION", DEFAULT_HELM_VERSION)
	var defaultHelmSprayVersion = getValueFromEnv("HELM_SPRAY_VERSION", DEFAULT_HELM_SPRAY_VERSION)
	initCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", defaultMinikubeUrl, "The URL to download minikube")
	initCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", defaultMinikubeVersion, "The minikube version")
	initCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", defaultDockerVersion, "The docker version")
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", defaultKubernetesVersion, "The kubernetes version")
	initCmd.Flags().StringVarP(&kubectlVersion, "kubectl-version", "", defaultKubectlVersion, "The kubectl version")
	initCmd.Flags().StringVarP(&helmVersion, "helm-version", "", defaultHelmVersion, "The helm version")
	initCmd.Flags().StringVarP(&helmSprayVersion, "helm-spray-version", "", defaultHelmSprayVersion, "The helm spray plugin version")
	initCmd.Flags().StringVarP(&sternVersion, "stern-version", "", DEFAULT_STERN_VERSION, "The stern version")
	initCmd.Flags().BoolVarP(&askForUpgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	initCmd.Flags().BoolVarP(&askForClean, "clean", "c", false, "Clean gokube (remove docker, minikube, kubectl and helm working directories)")
	initCmd.Flags().Int16VarP(&memory, "memory", "", int16(8192), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&cpus, "cpus", "", int16(4), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().StringVarP(&disk, "disk", "", "20g", "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().StringVarP(&checkIP, "check-ip", "", "192.168.99.100", "Checks if minikube VM allocated IP matches the provided one (0.0.0.0 means no check)")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", os.Getenv("INSECURE_REGISTRY"), "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", os.Getenv("HTTP_PROXY"), "HTTP proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", os.Getenv("HTTPS_PROXY"), "HTTPS proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", os.Getenv("NO_PROXY"), "No proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&miniappsRepo, "miniapps-repo", "", "https://thalesgroup.github.io/miniapps", "Helm repository for miniapps")
	initCmd.Flags().BoolVarP(&dnsProxy, "dns-proxy", "", false, "Use Virtualbox NAT DNS proxy (could be instable)")
	initCmd.Flags().BoolVarP(&hostDNSResolver, "host-dns-resolver", "", false, "Use Virtualbox NAT DNS host resolver (could be instable)")
	initCmd.Flags().BoolVarP(&debug, "debug", "", false, "Activate debug logging")
	initCmd.Flags().BoolVarP(&forceInit, "force", "f", false, "Force VM init (don't display warning message before initializing)")
	RootCmd.AddCommand(initCmd)
}

func getValueFromEnv(envVar string, defaultValue string) string {
	var value = os.Getenv(envVar)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}

func checkMinimumRequirements() {
	// Check minimum requirements
	if semver.New(minikubeVersion[1:]).Compare(*semver.New("1.6.1")) < 0 {
		fmt.Println("gokube is only compatible with minikube version >= 1.6.1")
		os.Exit(1)
	}
	if semver.New(helmVersion[1:]).Compare(*semver.New("2.16.1")) < 0 {
		fmt.Println("gokube is only compatible with helm version >= 2.16.1")
		os.Exit(1)
	}
}

func confirmCommandExecution() {
	fmt.Println("WARNING: Your Virtualbox GUI shall not be open and no other VM shall be currently running")
	fmt.Print("Press <CTRL+C> within the next 10s it you need to check this or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go gokube.WaitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
}

func resetVBLease() {
	// VB6 persists DHCP leases which prevent minikube to get the expected 192.168.99.100 IP address
	// Wait 5 seconds to make sure DHCP leases files are unlocked following VM deletion
	// TODO add manifest to ask for admin rights (when we will need to remove host-only network)
	fmt.Print("Resetting host-only network used by minikube...")
	var err error
	for n := 1; n < 3; n++ {
		time.Sleep(5 * time.Second)
		err = virtualbox.ResetHostOnlyNetworkLeases("192.168.99.1/24", debug)
		if err == nil {
			break
		} else {
			fmt.Print(".")
		}
	}
	if err != nil {
		fmt.Printf("\nCannot reset host-only network: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println()
	}
}

func exposeDashboard(port int) {
	for n := 1; n < 12; n++ {
		var dashboardService = kubectl.GetObject("kubernetes-dashboard", "svc", "kubernetes-dashboard")
		if len(dashboardService) > 0 {
			fmt.Println()
			patchPayload := fmt.Sprintf("{\"spec\":{\"type\":\"NodePort\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":9090,\"nodePort\":%d}]}}", port)
			kubectl.Patch("kubernetes-dashboard", "svc", "kubernetes-dashboard", patchPayload)
			break
		} else {
			fmt.Print(".")
			time.Sleep(5 * time.Second)
		}
	}
}

func installHelm() {
	// Initialize helm
	if !helm3 {
		helm.Init()
	}
	// Add helm miniapps repository
	helm.RepoAdd("miniapps", miniappsRepo)
	helm.RepoUpdate()
}

func installHelmPlugins() {
	if !helm3 {
		// Add helm spray plugin
		helmspray.DeletePlugin()
		helmspray.InstallPlugin(helmSprayVersion)
	} else {
		// TODO rely on helm plugin install
		fmt.Println("WARNING: helm-spray NOT installed as plugin installation is not yet compatible with helm3")
	}
}

func checkMinikubeIP() {
	var minikubeIP = minikube.Ip()
	if strings.Compare(checkIP, minikubeIP) != 0 {
		fmt.Printf("\nERROR: minikube IP (%s) does not match expected IP (%s), VM post-installation process aborted\n", minikubeIP, checkIP)
		os.Exit(1)
	}
}

func clean() {
	minikube.DeleteWorkingDirectory()
	helm.DeleteWorkingDirectory()
	kubectl.DeleteWorkingDirectory()
	docker.DeleteWorkingDirectory()
	docker.InitWorkingDirectory()
}

func upgrade() {
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

func initRun(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fmt.Println("usage: gokube init")
		os.Exit(1)
	}

	checkMinimumRequirements()

	helm3 = false
	if semver.New(helmVersion[1:]).Compare(*semver.New("3.0.0")) >= 0 {
		helm3 = true
	}
	ipCheckNeeded = strings.Compare("0.0.0.0", checkIP) != 0

	// Warn user with pre-requisites
	if ipCheckNeeded && !forceInit {
		confirmCommandExecution()
	}

	fmt.Println("Deleting previous minikube VM...")
	minikube.Delete()

	if ipCheckNeeded {
		resetVBLease()
	}

	// Force clean & upgrade if persisted gokube-version is lower than the current one
	if semver.New(gokubeVersion).Compare(*semver.New(GOKUBE_VERSION)) < 0 {
		fmt.Println("This version of gokube is launched for the first time, forcing clean & upgrade...")
		gokubeVersion = GOKUBE_VERSION
		askForClean = true
		askForUpgrade = true
	}

	if askForClean {
		fmt.Println("Deleting gokube dependencies working directory...")
		clean()
	}

	if askForUpgrade {
		fmt.Println("Downloading gokube dependencies...")
		upgrade()
	}

	// Keep kubernetes version in a persistent file to remember the right kubernetes version to set for (re)start command
	gokube.WriteConfig(gokubeVersion, kubernetesVersion)

	// Create virtual machine (minikube)
	fmt.Printf("Creating minikube VM with kubernetes %s...\n", kubernetesVersion)
	minikube.Start(memory, cpus, disk, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion, true, dnsProxy, hostDNSResolver)

	// Disable notification for updates
	minikube.ConfigSet("WantUpdateNotification", "false")
	// Enable dashboard
	minikube.AddonsEnable("dashboard")

	// Checks minikube IP
	if ipCheckNeeded {
		checkMinikubeIP()
	}

	// Switch context to minikube for kubectl and helm
	kubectl.ConfigUseContext("minikube")

	// Install helm
	fmt.Println("Initializing helm...")
	installHelm()
	if askForUpgrade {
		fmt.Println("Installing helm plugins...")
		installHelmPlugins()
	}

	// Patch kubernetes-dashboard to expose it on nodePort 30000
	fmt.Print("Exposing kubernetes dashboard to nodeport 30000...")
	exposeDashboard(30000)

	fmt.Println("\ngokube has been installed.")
}
