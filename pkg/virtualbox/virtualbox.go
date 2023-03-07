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

package virtualbox

import (
	"errors"
	"fmt"
	"github.com/gemalto/gokube/pkg/utils"
	"golang.org/x/sys/windows/registry"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	buggyNetmask = "0f000000"
)

// Host-only network.
type hostOnlyNetwork struct {
	Name        string
	GUID        string
	DHCP        bool
	IPv4        net.IPNet
	HwAddr      net.HardwareAddr
	Medium      string
	Status      string
	NetworkName string // referenced in DHCP.NetworkName
}

var (
	ErrNetworkAddrCidr = errors.New("host-only CIDR must be specified with a host address, not a network address")

	vboxManager = NewVBoxManager()
)

func IsRunning() (bool, error) {
	info, err := vboxManager.vbmOut("showvminfo", "minikube")
	if err != nil {
		return false, fmt.Errorf("not able to get VM info: %w", err)
	}
	re := regexp.MustCompile("(?m)^State:( *)(.*)$")
	rem := re.FindStringSubmatch(info)
	if len(rem) > 2 {
		if strings.HasPrefix(rem[2], "running") {
			return true, nil
		}
	}
	return false, nil
}

func Pause() error {
	err := vboxManager.vbm("controlvm", "minikube", "pause")
	if err != nil {
		return fmt.Errorf("not able to pause VM: %w", err)
	}
	return nil
}

func Resume() error {
	err := vboxManager.vbm("controlvm", "minikube", "resume")
	if err != nil {
		return fmt.Errorf("not able to resume VM: %w", err)
	}
	return nil
}

func Update(args ...string) error {
	err := vboxManager.vbm(append([]string{"modifyvm", "minikube"}, args...)...)
	if err != nil {
		return fmt.Errorf("not able to update VM: %w", err)
	}
	return nil
}

func DeleteSnapshot(name string) error {
	_, stderr, err := vboxManager.vbmOutErr("snapshot", "minikube", "delete", name)
	if err != nil {
		if reSnapshotNotFound.FindString(stderr) != "" || reNoSnapshotFound.FindString(stderr) != "" {
			return ErrSnapshotNotExist
		} else {
			return fmt.Errorf("not able to delete VM snapshot: %w", err)
		}
	}
	return nil
}

func TakeSnapshot(name string) error {
	err := vboxManager.vbm("snapshot", "minikube", "take", name)
	if err != nil {
		return fmt.Errorf("not able to take VM snapshot: %w", err)
	}
	return nil
}

func RestoreSnapshot(name string) error {
	err := vboxManager.vbm("snapshot", "minikube", "restore", name)
	if err != nil {
		return fmt.Errorf("not able to restore VM snapshot: %w", err)
	}
	return nil
}

func ResetHostOnlyNetworkLeases(hostOnlyCIDR string, verbose bool) error {
	nets, err := listHostOnlyAdapters(vboxManager)
	if err != nil {
		return fmt.Errorf("not able to list host-only network interfaces: %w", err)
	}
	if verbose {
		for _, v := range nets {
			fmt.Printf("\nResetHostOnlyNetworkLeases: listHostOnlyAdapters: %+v", v)
		}
	}
	ip, network, err := parseAndValidateCIDR(hostOnlyCIDR)
	if err != nil {
		return fmt.Errorf("not able to parse CIDR to find host-only network interface: %w", err)
	}
	if verbose {
		fmt.Printf("\nResetHostOnlyNetworkLeases: parseAndValidateCIDR: %s,%s", ip.String(), network.String())
	}
	hostOnlyNet := getHostOnlyAdapter(nets, ip, network.Mask)
	if hostOnlyNet == nil {
		if verbose {
			fmt.Printf("\nResetHostOnlyNetworkLeases: getHostOnlyAdapter: no host-only network interface matching minikube CDR")
		}
		return nil
	}
	if verbose {
		fmt.Printf("\nResetHostOnlyNetworkLeases: getHostOnlyAdapter: %s", hostOnlyNet.NetworkName)
	}
	filesPattern := utils.GetUserHome() + "/.VirtualBox/" + hostOnlyNet.NetworkName + "*"
	files, err := filepath.Glob(filesPattern)
	if err != nil {
		return fmt.Errorf("not able to get host-only network interface DHCP leases files: %w", err)
	}
	for _, f := range files {
		if verbose {
			fmt.Printf("\nResetHostOnlyNetworkLeases: deleting lease file %s...", f)
		}
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("not able to delete lease file %s: %w", f, err)
		}
		if verbose {
			fmt.Printf("\nResetHostOnlyNetworkLeases: deleted lease file %s", f)
		}
	}
	return nil
}

func listHostOnlyAdapters(vbox VBoxManager) (map[string]*hostOnlyNetwork, error) {
	out, err := vbox.vbmOut("list", "hostonlyifs")
	if err != nil {
		return nil, err
	}

	byName := map[string]*hostOnlyNetwork{}
	byIP := map[string]*hostOnlyNetwork{}
	n := &hostOnlyNetwork{}

	err = parseKeyValues(out, reColonLine, func(key, val string) error {
		switch key {
		case "Name":
			n.Name = val
		case "GUID":
			n.GUID = val
		case "DHCP":
			n.DHCP = (val != "Disabled")
		case "IPAddress":
			n.IPv4.IP = net.ParseIP(val)
		case "NetworkMask":
			n.IPv4.Mask = parseIPv4Mask(val)
		case "HardwareAddress":
			mac, err := net.ParseMAC(val)
			if err != nil {
				return err
			}
			n.HwAddr = mac
		case "MediumType":
			n.Medium = val
		case "Status":
			n.Status = val
		case "VBoxNetworkName":
			n.NetworkName = val

			if _, present := byName[n.NetworkName]; present {
				return fmt.Errorf("VirtualBox is configured with multiple host-only adapters with the same name %q. Please remove one", n.NetworkName)
			}
			byName[n.NetworkName] = n

			if len(n.IPv4.IP) != 0 {
				if _, present := byIP[n.IPv4.IP.String()]; present {
					return fmt.Errorf("VirtualBox is configured with multiple host-only adapters with the same IP %q. Please remove one", n.IPv4.IP)
				}
				byIP[n.IPv4.IP.String()] = n
			}

			n = &hostOnlyNetwork{}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return byName, nil
}

func getHostOnlyAdapter(nets map[string]*hostOnlyNetwork, hostIP net.IP, netmask net.IPMask) *hostOnlyNetwork {
	for _, n := range nets {
		// Second part of this conditional handles a race where
		// VirtualBox returns us the incorrect netmask value for the
		// newly created adapter.
		if hostIP.Equal(n.IPv4.IP) &&
			(netmask.String() == n.IPv4.Mask.String() || n.IPv4.Mask.String() == buggyNetmask) {
			return n
		}
	}

	return nil
}

func parseAndValidateCIDR(hostOnlyCIDR string) (net.IP, *net.IPNet, error) {
	ip, network, err := net.ParseCIDR(hostOnlyCIDR)
	if err != nil {
		return nil, nil, err
	}

	networkAddress := network.IP.To4()
	if ip.Equal(networkAddress) {
		return nil, nil, ErrNetworkAddrCidr
	}

	return ip, network, nil
}

func detectVBoxManageCmd() string {
	cmd := "VBoxManage"
	if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" {
		if path, err := exec.LookPath(filepath.Join(p, cmd)); err == nil {
			return path
		}
	}

	if p := os.Getenv("VBOX_MSI_INSTALL_PATH"); p != "" {
		if path, err := exec.LookPath(filepath.Join(p, cmd)); err == nil {
			return path
		}
	}

	// Look in default installation path for VirtualBox version > 5
	if path, err := exec.LookPath(filepath.Join("C:\\Program Files\\Oracle\\VirtualBox", cmd)); err == nil {
		return path
	}

	// Look in windows registry
	if p, err := findVBoxInstallDirInRegistry(); err == nil {
		if path, err := exec.LookPath(filepath.Join(p, cmd)); err == nil {
			return path
		}
	}

	return detectVBoxManageCmdInPath() //fallback to path
}

func detectVBoxManageCmdInPath() string {
	cmd := "VBoxManage"
	if path, err := exec.LookPath(cmd); err == nil {
		return path
	}
	return cmd
}

func findVBoxInstallDirInRegistry() (string, error) {
	registryKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Oracle\VirtualBox`, registry.QUERY_VALUE)
	if err != nil {
		errorMessage := fmt.Sprintf("Can't find VirtualBox registry entries, is VirtualBox really installed properly? %s", err)
		return "", fmt.Errorf(errorMessage)
	}

	defer registryKey.Close()

	installDir, _, err := registryKey.GetStringValue("InstallDir")
	if err != nil {
		errorMessage := fmt.Sprintf("Can't find InstallDir registry key within VirtualBox registries entries, is VirtualBox really installed properly? %s", err)
		return "", fmt.Errorf(errorMessage)
	}

	return installDir, nil
}

func parseIPv4Mask(s string) net.IPMask {
	mask := net.ParseIP(s)
	if mask == nil {
		return nil
	}
	return net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
}
