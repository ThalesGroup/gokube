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

// TODO manage vbox time sync (VBoxManage guestproperty set default "/VirtualBox/GuestAdd/VBoxService/--timesync-set-threshold" 1000)

import (
	"fmt"
	"github.com/gemalto/gokube/internal/util"
	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
	"os/exec"
)

var memory int16
var cpus int16
var swap int16
var enableSwap bool
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
var keepVM bool
var dnsDomain string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern + k9s and creates a minikube VM",
	Long:         "Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern + k9s and creates a minikube VM",
	RunE:         initRun,
	SilenceUsage: true,
}

func init() {
	defaultVMMemory, _ := strconv.Atoi(utils.GetValueFromEnv("MINIKUBE_MEMORY", strconv.Itoa(DEFAULT_MINIKUBE_MEMORY)))
	defaultVMCPUs, _ := strconv.Atoi(utils.GetValueFromEnv("MINIKUBE_CPUS", strconv.Itoa(DEFAULT_MINIKUBE_CPUS)))
	defaultVMSwap, _ := strconv.Atoi(utils.GetValueFromEnv("MINIKUBE_SWAP", strconv.Itoa(DEFAULT_MINIKUBE_SWAP)))
	enableSwap = false
	if defaultVMSwap != 0 {
		enableSwap = true
	}
	defaultGokubeQuiet := false
	if len(utils.GetValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	loadURLVersionsFromEnv()
	initCmd.Flags().StringVarP(&kubernetesVersion, "kubernetes-version", "", utils.GetValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBERNETES_VERSION), "The kubernetes version")
	initCmd.Flags().StringVarP(&containerRuntime, "container-runtime", "", utils.GetValueFromEnv("MINIKUBE_CONTAINER_RUNTIME", DEFAULT_MINIKUBE_CONTAINER_RUNTIME), "Minikube container runtime (docker, cri-o, containerd)")
	initCmd.Flags().BoolVarP(&askForUpgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	initCmd.Flags().BoolVarP(&askForClean, "clean", "c", false, "Clean gokube (remove docker, minikube, kubectl and helm working directories)")
	initCmd.Flags().Int16VarP(&memory, "memory", "", int16(defaultVMMemory), "Amount of RAM allocated to the minikube VM in MB")
	initCmd.Flags().Int16VarP(&cpus, "cpus", "", int16(defaultVMCPUs), "Number of CPUs allocated to the minikube VM")
	initCmd.Flags().Int16VarP(&swap, "swap", "", int16(defaultVMSwap), "Amount of SWAP allocated to the minikube VM in MB")
	initCmd.Flags().StringVarP(&disk, "disk", "", utils.GetValueFromEnv("MINIKUBE_DISK", DEFAULT_MINIKUBE_DISK), "Disk size allocated to the minikube VM. Format: <number>[<unit>], where unit = b, k, m or g")
	initCmd.Flags().StringVarP(&checkIP, "check-ip", "", utils.GetValueFromEnv("GOKUBE_CHECK_IP", DEFAULT_GOKUBE_CHECK_IP), "Checks if minikube VM allocated IP matches the provided one (0.0.0.0 means no check)")
	initCmd.Flags().StringVarP(&insecureRegistry, "insecure-registry", "", os.Getenv("INSECURE_REGISTRY"), "Insecure Docker registries to pass to the Docker daemon. The default service CIDR range will automatically be added.")
	initCmd.Flags().StringVarP(&httpProxy, "http-proxy", "", os.Getenv("HTTP_PROXY"), "HTTP proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&httpsProxy, "https-proxy", "", os.Getenv("HTTPS_PROXY"), "HTTPS proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&noProxy, "no-proxy", "", os.Getenv("NO_PROXY"), "No proxy variable for docker engine in minikube VM")
	initCmd.Flags().StringVarP(&dnsDomain, "dns-domain", "", utils.GetValueFromEnv("MINIKUBE_DNS_DOMAIN", DEFAULT_MINIKUBE_DNS_DOMAIN), "Minikube cluster DNS domain name")
	initCmd.Flags().BoolVarP(&dnsProxy, "dns-proxy", "", false, "Use Virtualbox NAT DNS proxy (could be unstable)")
	initCmd.Flags().BoolVarP(&hostDNSResolver, "host-dns-resolver", "", false, "Use Virtualbox NAT DNS host resolver (could be unstable)")
	initCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before initializing")
	initCmd.Flags().BoolVar(&keepVM, "keep-vm", false, "Keep minikube VM as it is (don't delete/recreate)")
	initCmd.Flags().BoolVar(&force, "force", false, "Force minikube to perform possibly dangerous operations")
	rootCmd.AddCommand(initCmd)
}

func resetVBLease(hostOnlyCIDR string) error {
	// VB6 persists DHCP leases which prevent minikube to get the expected 192.168.99.100 IP address
	// Wait 5 seconds to make sure DHCP leases files are unlocked following VM deletion
	// TODO add manifest to ask for admin rights (when we will need to remove host-only network)
	fmt.Print("Resetting host-only network used by minikube...")
	var err error
	for n := 1; n < 3; n++ {
		time.Sleep(5 * time.Second)
		err = virtualbox.ResetHostOnlyNetworkLeases(hostOnlyCIDR, verbose)
		if err == nil {
			break
		} else {
			fmt.Print(".")
		}
	}
	if err != nil {
		return fmt.Errorf("cannot reset host-only network: %w", err)
	} else {
		fmt.Println()
		return nil
	}
}

func setupMiniappsHelmRepository() error {
	err := helm.RepoAdd("miniapps", miniappsRepo)
	if err != nil {
		return fmt.Errorf("cannot add miniapps repo: %w", err)
	}
	err = helm.RepoUpdate()
	if err != nil {
		return fmt.Errorf("cannot update helm repositories: %w", err)
	}
	return nil
}

func installChartMuseum(localRepoIp string) error {
	err := helm.RepoAdd("chartmuseum", "https://chartmuseum.github.io/charts")
	if err != nil {
		return fmt.Errorf("cannot add chartmuseum repo: %w", err)
	}
	err = helm.RepoUpdate()
	if err != nil {
		return fmt.Errorf("cannot update helm repositories: %w", err)
	}
	err = helm.Upgrade("chartmuseum/chartmuseum", "", "chartmuseum", "kube-system", "env.open.DISABLE_API=false,env.open.ALLOW_OVERWRITE=true,service.type=NodePort,service.nodePort=32767", "")
	if err != nil {
		return fmt.Errorf("cannot install chartmuseum: %w", err)
	}
	fmt.Printf("Waiting for chartmuseum...")
	retries := 18
	waitBeforeRetry := 5
	for n := 1; n <= retries; n++ {
		readyReplicas, err := kubectl.Get("kube-system", "deploy", "chartmuseum", "{.status.readyReplicas}")
		if err != nil {
			return fmt.Errorf("cannot get K8S chartmuseum deployment: %w", err)
		}
		var ready int
		if len(readyReplicas) > 0 {
			n, err := strconv.Atoi(readyReplicas)
			if err != nil {
				return fmt.Errorf("cannot check chartmuseum readiness: %w", err)
			}
			ready = n
		}
		if ready > 0 {
			fmt.Println()
			break
		} else {
			fmt.Print(".")
			if n == retries {
				fmt.Printf("\nWarning: chartmuseum is not ready after %ds, which probably means its installation failed\n", retries*waitBeforeRetry)
			} else {
				time.Sleep(time.Duration(waitBeforeRetry) * time.Second)
			}
		}
	}
	err = helm.RepoAdd("minikube", "http://"+localRepoIp+":32767")
	if err != nil {
		fmt.Printf("Warning: cannot add minikube repo: %s\n", err)
	}
	err = helm.RepoUpdate()
	if err != nil {
		return fmt.Errorf("cannot update helm repositories: %w", err)
	}
	return nil
}

func exposeDashboard(port int) error {
	for n := 1; n <= 12; n++ {
		dashboardService, err := kubectl.Get("kubernetes-dashboard", "svc", "kubernetes-dashboard", "")
		if err != nil {
			return fmt.Errorf("cannot get K8S kubernetes-dashboard service: %w", err)
		}
		if len(dashboardService) > 0 {
			fmt.Println()
			patchPayload := fmt.Sprintf("{\"spec\":{\"type\":\"NodePort\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":9090,\"nodePort\":%d}]}}", port)
			err = kubectl.Patch("kubernetes-dashboard", "svc", "kubernetes-dashboard", patchPayload)
			if err != nil {
				return fmt.Errorf("cannot patch K8S kubernetes-dashboard service: %w", err)
			}
			break
		} else {
			fmt.Print(".")
			if n == 12 {
				fmt.Printf("\nWarning: kubernetes-dashboard is not present after 60s, which probably means its installation failed\n")
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}
	return nil
}

func initRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}

	checkLatestVersion()

	err := gokube.ReadConfig(verbose)
	if err != nil {
		return fmt.Errorf("cannot read gokube configuration file: %w", err)
	}
	gokubeVersion = viper.GetString("gokube-version")
	if len(gokubeVersion) == 0 {
		gokubeVersion = "0.0.0"
	}

	// Force clean & upgrade if persisted gokube-version is lower than the current one
	if semver.New(gokubeVersion).Compare(*semver.New(GOKUBE_VERSION)) < 0 {
		fmt.Println("Warning: this version of gokube is launched for the first time, forcing clean & upgrade...")
		gokubeVersion = GOKUBE_VERSION
		askForClean = true
		askForUpgrade = true
	}

	if !askForUpgrade {
		checkMinimumRequirements()
	}

	ipCheckNeeded = strings.Compare("0.0.0.0", checkIP) != 0

	if askForClean && keepVM {
		fmt.Println("Error: Cannot keep VM while cleaning gokube")
		os.Exit(1)
	}

	// Warn user with pre-requisites
	if ipCheckNeeded && !quiet && !keepVM {
		gokube.ConfirmInitCommandExecution()
	}

	startTime := time.Now()

	if !keepVM {
		fmt.Println("Deleting previous minikube VM...")
		err := minikube.Delete()
		if err != nil {
			fmt.Printf("Warning: cannot delete previous minikube VM: %s\n", err)
		}
		if ipCheckNeeded {
			err = resetVBLease(DEFAULT_GOKUBE_CIDR)
			if err != nil {
				return fmt.Errorf("cannot delete previous minikube VM: %w", err)
			}
		}
	}

	if askForClean {
		fmt.Println("Deleting gokube dependencies working directory...")
		_ = minikube.DeleteWorkingDirectory()
		_ = kubectl.DeleteWorkingDirectory()
		_ = docker.DeleteWorkingDirectory()
		_ = docker.InitWorkingDirectory()
		_ = helm.DeleteWorkingDirectory()
	} else if !keepVM {
		_ = helm.ResetWorkingDirectory()
	}

	if askForUpgrade {
		fmt.Println("Upgrading gokube dependencies...")
		err = upgradeDependencies()
		if err != nil {
			return err
		}
	}

	if !keepVM {
		// Disable notification for updates
		_ = minikube.ConfigSet("WantUpdateNotification", "false")

		// Create virtual machine (minikube)
		fmt.Printf("Creating minikube VM with kubernetes %s...\n", kubernetesVersion)
		err := minikube.Start(memory, cpus, disk, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion, true, dnsProxy, hostDNSResolver, dnsDomain, containerRuntime, force, verbose)
		if err != nil {
			return fmt.Errorf("cannot start minikube VM: %w", err)
		}

        // Create & attach swap drive to minikube
        if enableSwap {
            fmt.Println("Creating & attaching swap drive to minikube VM...")
            vboxManager := virtualbox.NewVBoxManager()
            err = vboxManager.AddSwapDisk(swap)
            if err != nil {
                fmt.Printf("Warning: cannot create & attach swap drive to minikube VM: %s\n", err)
            }
        }

		// Enable dashboard
		err = minikube.AddonsEnable("dashboard")
		if err != nil {
			return fmt.Errorf("cannot enable dashboard minikube add-on: %w", err)
		}

		minikubeIP, err := minikube.Ip()
		if err != nil {
			return fmt.Errorf("cannot get minikube VM IP address: %w", err)
		}
		// Checks minikube IP
		if ipCheckNeeded {
			if strings.Compare(checkIP, minikubeIP) != 0 {
				fmt.Printf("\nError: minikube IP (%s) does not match expected IP (%s), VM post-installation process aborted\n", minikubeIP, checkIP)
				os.Exit(1)
			}
		}

		// Switch context to minikube for kubectl and helm
		err = kubectl.ConfigUseContext("minikube")
		if err != nil {
			return fmt.Errorf("cannot switch K8S context to minikube: %w", err)
		}

		fmt.Println("Installing ChartMuseum...")
		err = installChartMuseum(minikubeIP)
		if err != nil {
			return err
		}

		fmt.Println("Configuring miniapps repository...")
		err = setupMiniappsHelmRepository()
		if err != nil {
			return err
		}

		// Patch kubernetes-dashboard to expose it on nodePort 30000
		fmt.Print("Exposing kubernetes dashboard to nodeport 30000...")
		err = exposeDashboard(30000)
		if err != nil {
			return err
		}

	}

	if askForUpgrade {
		fmt.Println("Upgrading helm plugins...")
		err := upgradeHelmPlugins()
		if err != nil {
			return err
		}
	}

	// Keep kubernetes version in a persistent file to remember the right kubernetes version to set for (re)start command
	err = gokube.WriteConfig(gokubeVersion, kubernetesVersion, containerRuntime)
	if err != nil {
		return fmt.Errorf("cannot write gokube configuration: %w", err)
	}

    // Format & enable swap drive in minikube VM
    if enableSwap {
        fmt.Println("Formatting & enabling swap drive in minikube VM...")
	    err = addSwapToMinikube()
	    if err != nil {
		    fmt.Printf("Warning: cannot format/enable swap drive in minikube VM: %s\n", err)
	    }
    }

	fmt.Printf("\ngokube init completed in %s\n", util.Duration(time.Since(startTime)))
	return nil
}

func addSwapToMinikube() error {

	// Add swap file commands
	swapCmds := []string{
		"sudo mkswap /dev/sdb",
		"sudo swapon /dev/sdb",
		"echo '/dev/sdb none swap defaults 0 0' | sudo tee -a /etc/fstab",
	}

	// Execute each command
	for _, cmd := range swapCmds {
		sshCmd := exec.Command("minikube", "ssh", cmd)
		err := sshCmd.Run()
		if err != nil {
			return fmt.Errorf("error running command '%s': %w", cmd, err)
		}
	}

	return nil
}