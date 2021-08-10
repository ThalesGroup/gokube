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
	"github.com/gemalto/gokube/internal/util"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/viper"
	"os"
	"strconv"
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

var kubernetesVersion string
var memory int16
var cpus int16
var disk string
var checkIP string
var ipCheckNeeded bool
var insecureRegistry string
var httpProxy string
var httpsProxy string
var noProxy string
var askForClean bool
var miniappsRepo string
var dnsProxy bool
var hostDNSResolver bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates a minikube VM",
	Long:         "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates a minikube VM",
	RunE:         initRun,
	SilenceUsage: true,
}

func init() {
	var defaultKubernetesVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBERNETES_VERSION)
	var defaultKubectlVersion = getValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBECTL_VERSION)
	var defaultMinikubeUrl = getValueFromEnv("MINIKUBE_URL", DEFAULT_MINIKUBE_URL)
	var defaultMinikubeVersion = getValueFromEnv("MINIKUBE_VERSION", DEFAULT_MINIKUBE_VERSION)
	var defaultDockerVersion = getValueFromEnv("DOCKER_VERSION", DEFAULT_DOCKER_VERSION)
	var defaultHelmVersion = getValueFromEnv("HELM_VERSION", DEFAULT_HELM_VERSION)
	var defaultHelmSprayUrl = getValueFromEnv("HELM_SPRAY_URL", DEFAULT_HELM_SPRAY_URL)
	var defaultHelmSprayVersion = getValueFromEnv("HELM_SPRAY_VERSION", DEFAULT_HELM_SPRAY_VERSION)
	var defaultHelmImageUrl = getValueFromEnv("HELM_IMAGE_URL", DEFAULT_HELM_IMAGE_URL)
	var defaultHelmImageVersion = getValueFromEnv("HELM_IMAGE_VERSION", DEFAULT_HELM_IMAGE_VERSION)
	defaultVMMemory, _ := strconv.Atoi(getValueFromEnv("MINIKUBE_MEMORY", strconv.Itoa(DEFAULT_MINIKUBE_MEMORY)))
	defaultVMCPUs, _ := strconv.Atoi(getValueFromEnv("MINIKUBE_CPUS", strconv.Itoa(DEFAULT_MINIKUBE_CPUS)))
	defaultGokubeQuiet := false
	if len(getValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	initCmd.Flags().StringVarP(&minikubeURL, "minikube-url", "", defaultMinikubeUrl, "The URL to download minikube")
	initCmd.Flags().StringVarP(&minikubeVersion, "minikube-version", "", defaultMinikubeVersion, "The minikube version")
	initCmd.Flags().StringVarP(&dockerVersion, "docker-version", "", defaultDockerVersion, "The docker version")
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", defaultKubernetesVersion, "The kubernetes version")
	initCmd.Flags().StringVarP(&kubectlVersion, "kubectl-version", "", defaultKubectlVersion, "The kubectl version")
	initCmd.Flags().StringVarP(&helmVersion, "helm-version", "", defaultHelmVersion, "The helm version")
	initCmd.Flags().StringVarP(&helmSprayURL, "helm-spray-url", "", defaultHelmSprayUrl, "The URL to download helm spray plugin")
	initCmd.Flags().StringVarP(&helmSprayVersion, "helm-spray-version", "", defaultHelmSprayVersion, "The helm spray plugin version")
	initCmd.Flags().StringVarP(&helmImageURL, "helm-image-url", "", defaultHelmImageUrl, "The URL to download helm image plugin")
	initCmd.Flags().StringVarP(&helmImageVersion, "helm-image-version", "", defaultHelmImageVersion, "The helm image image version")
	initCmd.Flags().StringVarP(&sternVersion, "stern-version", "", DEFAULT_STERN_VERSION, "The stern version")
	initCmd.Flags().BoolVarP(&askForUpgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	initCmd.Flags().BoolVarP(&askForClean, "clean", "c", false, "Clean gokube (remove docker, minikube, kubectl and helm working directories)")
	initCmd.Flags().Int16VarP(&memory, "memory", "", int16(defaultVMMemory), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&cpus, "cpus", "", int16(defaultVMCPUs), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().StringVarP(&disk, "disk", "", "20g", "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().StringVarP(&checkIP, "check-ip", "", "192.168.99.100", "Checks if minikube VM allocated IP matches the provided one (0.0.0.0 means no check)")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", os.Getenv("INSECURE_REGISTRY"), "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", os.Getenv("HTTP_PROXY"), "HTTP proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", os.Getenv("HTTPS_PROXY"), "HTTPS proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", os.Getenv("NO_PROXY"), "No proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&miniappsRepo, "miniapps-repo", "", DEFAULT_MINIAPPS_REPO, "Helm repository for miniapps")
	initCmd.Flags().BoolVarP(&dnsProxy, "dns-proxy", "", false, "Use Virtualbox NAT DNS proxy (could be instable)")
	initCmd.Flags().BoolVarP(&hostDNSResolver, "host-dns-resolver", "", false, "Use Virtualbox NAT DNS host resolver (could be instable)")
	initCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before initializing")
	rootCmd.AddCommand(initCmd)
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
	if semver.New(kubernetesVersion[1:]).Compare(*semver.New("1.8.0")) < 0 {
		fmt.Println("FATAL: This gokube version is only compatible with kubernetes version >= 1.8.0")
		os.Exit(1)
	}
	if semver.New(minikubeVersion[1:]).Compare(*semver.New("1.6.1")) < 0 {
		fmt.Println("FATAL: This gokube version is only compatible with minikube version >= 1.6.1")
		os.Exit(1)
	}
	if semver.New(helmVersion[1:]).Compare(*semver.New("3.0.0-0")) < 0 {
		fmt.Println("FATAL: This gokube version is only compatible with helm version >= 3.0.0-0")
		os.Exit(1)
	}
	if semver.New(helmSprayVersion[1:]).Compare(*semver.New("4.0.0-0")) < 0 {
		fmt.Println("FATAL: This gokube version is only compatible with helm-spray version >= 4.0.0-0")
		os.Exit(1)
	}
}

func confirmInitCommandExecution() {
	fmt.Println("WARNING: Your Virtualbox GUI shall not be open and no other VM shall be currently running")
	fmt.Print("Press <CTRL+C> within the next 10s it you need to check this or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go gokube.WaitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func resetVBLease() {
	// VB6 persists DHCP leases which prevent minikube to get the expected 192.168.99.100 IP address
	// Wait 5 seconds to make sure DHCP leases files are unlocked following VM deletion
	// TODO add manifest to ask for admin rights (when we will need to remove host-only network)
	fmt.Print("Resetting host-only network used by minikube...")
	var err error
	for n := 1; n < 3; n++ {
		time.Sleep(5 * time.Second)
		err = virtualbox.ResetHostOnlyNetworkLeases("192.168.99.1/24", verbose)
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
	for n := 1; n <= 12; n++ {
		var dashboardService = kubectl.Get("kubernetes-dashboard", "svc", "kubernetes-dashboard", "")
		if len(dashboardService) > 0 {
			fmt.Println()
			patchPayload := fmt.Sprintf("{\"spec\":{\"type\":\"NodePort\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":9090,\"nodePort\":%d}]}}", port)
			kubectl.Patch("kubernetes-dashboard", "svc", "kubernetes-dashboard", patchPayload)
			break
		} else {
			fmt.Print(".")
			if n == 12 {
				fmt.Printf("\nWARNING: kubernetes-dashboard is not present after 60s, which probably means its installation failed\n")
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func waitChartMuseum() {
	retries := 18
	for n := 1; n <= retries; n++ {
		readyReplicas := kubectl.Get("kube-system", "deploy", "chartmuseum", "{.status.readyReplicas}")
		var ready int
		if len(readyReplicas) > 0 {
			n, err := strconv.Atoi(readyReplicas)
			if err != nil {
				fmt.Printf("\nCannot check chartmuseum readiness: %s\n", err)
				os.Exit(1)
			}
			ready = n
		}
		if ready > 0 {
			fmt.Println()
			break
		} else {
			fmt.Print(".")
			if n == retries {
				fmt.Printf("\nWARNING: chartmuseum is not ready after 90s, which probably means its installation failed\n")
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func configureHelm(localRepoIp string) {
	// Add helm chartmuseum and miniapps repository
	helm.RepoAdd("chartmuseum", "https://chartmuseum.github.io/charts")
	helm.RepoAdd("miniapps", miniappsRepo)
	helm.RepoUpdate()
	// Install chartmuseum
	helm.Upgrade("chartmuseum/chartmuseum", "", "chartmuseum", "kube-system", "env.open.DISABLE_API=false,env.open.ALLOW_OVERWRITE=true,service.type=NodePort,service.nodePort=32767", "")
	fmt.Printf("Waiting for chartmuseum...")
	waitChartMuseum()
	helm.RepoAdd("minikube", "http://"+localRepoIp+":32767")
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
	kubectl.DeleteWorkingDirectory()
	docker.DeleteWorkingDirectory()
	docker.InitWorkingDirectory()
	helm.DeleteWorkingDirectory()
}

// TODO manage vbox time sync
// VBoxManage guestproperty set default "/VirtualBox/GuestAdd/VBoxService/--timesync-set-threshold" 1000

func initRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}

	checkMinimumRequirements()
	checkLatestVersion()

	ipCheckNeeded = strings.Compare("0.0.0.0", checkIP) != 0

	// Warn user with pre-requisites
	if ipCheckNeeded && !quiet {
		confirmInitCommandExecution()
	}

	startTime := time.Now()

	fmt.Println("Deleting previous minikube VM...")
	minikube.Delete()

	if ipCheckNeeded {
		resetVBLease()
	}

	gokube.ReadConfig(verbose)
	gokubeVersion = viper.GetString("gokube-version")
	if len(gokubeVersion) == 0 {
		gokubeVersion = "0.0.0"
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
	helm.ResetWorkingDirectory()

	if askForUpgrade {
		fmt.Println("Downloading gokube dependencies...")
		upgrade()
	}

	// Disable notification for updates
	minikube.ConfigSet("WantUpdateNotification", "false")

	// Create virtual machine (minikube)
	fmt.Printf("Creating minikube VM with kubernetes %s...\n", kubernetesVersion)
	minikube.Start(memory, cpus, disk, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion, true, dnsProxy, hostDNSResolver)

	// Enable dashboard
	minikube.AddonsEnable("dashboard")

	// Checks minikube IP
	if ipCheckNeeded {
		checkMinikubeIP()
	}

	// Switch context to minikube for kubectl and helm
	kubectl.ConfigUseContext("minikube")

	// Install helm
	fmt.Println("Configuring helm...")
	configureHelm(minikube.Ip())
	if askForUpgrade {
		fmt.Println("Installing helm plugins...")
		installHelmPlugins()
	}

	// Patch kubernetes-dashboard to expose it on nodePort 30000
	fmt.Print("Exposing kubernetes dashboard to nodeport 30000...")
	exposeDashboard(30000)

	// Keep kubernetes version in a persistent file to remember the right kubernetes version to set for (re)start command
	gokube.WriteConfig(gokubeVersion, kubernetesVersion)

	fmt.Printf("\ngokube setup completed in %s\n", util.Duration(time.Since(startTime)))
	return nil
}
