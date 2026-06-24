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
	"github.com/gemalto/gokube/pkg/hypervisor"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
	"net/http"

	"github.com/coreos/go-semver/semver"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
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
	initCmd.Flags().StringVarP(&driver, "driver", "", utils.GetValueFromEnv("MINIKUBE_DRIVER", DEFAULT_MINIKUBE_DRIVER), "Minikube driver/hypervisor (virtualbox, hyperv)")
	initCmd.Flags().StringVarP(&hypervVirtualSwitch, "hyperv-virtual-switch", "", utils.GetValueFromEnv("MINIKUBE_HYPERV_VIRTUAL_SWITCH", ""), "Name of the Hyper-V virtual switch to use (driver=hyperv only; empty lets minikube auto-select the Default Switch)")
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

	// Install or Upgrade the ChartMuseum
	err = helm.Upgrade("chartmuseum/chartmuseum", "", "chartmuseum", "kube-system",
		"env.open.DISABLE_API=false,env.open.ALLOW_OVERWRITE=true,service.type=NodePort,service.nodePort=32767", "")
	if err != nil {
		return fmt.Errorf("cannot install chartmuseum: %w", err)
	}

	fmt.Printf("Waiting for chartmuseum to be ready...\n")

	retries := 50
	waitBeforeRetry := 5

	for n := 1; n <= retries; n++ {
		// Check if ChartMuseum service is ready
		ready, err := isChartMuseumReady(localRepoIp, 32767)
		if err == nil && ready {
			fmt.Println("ChartMuseum is ready.")
			break
		} else {
			fmt.Print(".")
			if n == retries {
				fmt.Printf("\nWarning: chartmuseum is not ready after %ds, which likely means its installation failed\n", retries*waitBeforeRetry)
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

// Function to check if ChartMuseum is ready
func isChartMuseumReady(ip string, port int) (bool, error) {
	url := fmt.Sprintf("http://%s:%d/index.yaml", ip, port)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
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

	checkMinimumRequirements()

	// Resolve the hypervisor for the requested driver and fail fast if the host
	// is not ready (e.g. Hyper-V not enabled, not running elevated, switch
	// missing) before performing any destructive VM operation.
	hv, err := hypervisor.New(driver)
	if err != nil {
		return fmt.Errorf("invalid minikube driver %q: %w", driver, err)
	}
	if err = hv.Validate(hypervVirtualSwitch); err != nil {
		return err
	}

	// Drivers other than virtualbox (hyperv) have no predictable static IP, so
	// disable the IP check. This single toggle also skips the VirtualBox-only
	// host-only lease reset and the pre-init VirtualBox GUI warning below.
	if driver != hypervisor.DriverVirtualBox {
		checkIP = "0.0.0.0"
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
		// Host-only DHCP lease reset is VirtualBox-specific.
		if ipCheckNeeded && driver == hypervisor.DriverVirtualBox {
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
		gokube.DeleteAllExecutables()
	} else if !keepVM {
		_ = helm.ResetWorkingDirectory()
	}

	fmt.Println("Checking gokube dependencies...")
	err = upgradeDependencies()
	if err != nil {
		return err
	}

	// Persist gokube-version immediately after downloads succeed.
	// If any subsequent step (VM creation, ChartMuseum, helm repos) fails and
	// initRun returns early, the version is already written. The next run will
	// read "1.38.0", skip the forced-clean path, and serve all tools from cache
	// instead of wiping metadata and re-downloading everything.
	if err = gokube.WriteConfig(gokubeVersion, kubernetesVersion, containerRuntime, driver, hypervVirtualSwitch); err != nil {
		return fmt.Errorf("cannot write gokube configuration: %w", err)
	}

	if !keepVM {
		// Disable notification for updates
		_ = minikube.ConfigSet("WantUpdateNotification", "false")

		// Create virtual machine (minikube)
		fmt.Printf("Creating minikube VM with kubernetes %s...\n", kubernetesVersion)
		err := minikube.Start(memory, cpus, disk, httpProxy, httpsProxy, noProxy, insecureRegistry, kubernetesVersion, true, dnsProxy, hostDNSResolver, dnsDomain, containerRuntime, force, verbose, driver, hypervVirtualSwitch)
		if err != nil {
			return fmt.Errorf("cannot start minikube VM: %w", err)
		}

        // Create, attach, detect & enable swap drive in minikube VM.
        // The in-VM device node is detected (not assumed to be /dev/sdb) so the
        // logic is correct for both VirtualBox and Hyper-V. Swap on Hyper-V is
        // experimental (see AddSwapDisk).
        if enableSwap {
            fmt.Println("Creating & attaching swap drive to minikube VM...")
            before, listErr := listMinikubeDisks()
            err = hv.AddSwapDisk(swap)
            if err != nil {
                fmt.Printf("Warning: cannot create & attach swap drive to minikube VM: %s\n", err)
            } else if listErr != nil {
                fmt.Printf("Warning: cannot list minikube disks to locate swap drive: %s\n", listErr)
            } else if swapDevice, derr := detectNewSwapDevice(before); derr != nil {
                fmt.Printf("Warning: cannot locate attached swap drive in minikube VM: %s\n", derr)
            } else {
                fmt.Println("Formatting & enabling swap drive in minikube VM...")
                if ferr := formatAndEnableSwap(swapDevice); ferr != nil {
                    fmt.Printf("Warning: cannot format/enable swap drive in minikube VM: %s\n", ferr)
                }
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

	fmt.Println("Checking helm plugins...")
	if err := upgradeHelmPlugins(); err != nil {
		return err
	}

	// Keep kubernetes version in a persistent file to remember the right kubernetes version to set for (re)start command
	err = gokube.WriteConfig(gokubeVersion, kubernetesVersion, containerRuntime, driver, hypervVirtualSwitch)
	if err != nil {
		return fmt.Errorf("cannot write gokube configuration: %w", err)
	}

	fmt.Printf("\ngokube init completed in %s\n", util.Duration(time.Since(startTime)))
	return nil
}