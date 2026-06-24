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

package hypervisor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gemalto/gokube/pkg/utils"
)

// vmName is the name minikube gives the Hyper-V VM (it matches the profile).
const vmName = "minikube"

// hypervHypervisor implements Hypervisor on top of the Hyper-V PowerShell
// module. Every operation shells out via runPowerShell.
type hypervHypervisor struct{}

func (h *hypervHypervisor) IsRunning() (bool, error) {
	// -ErrorAction SilentlyContinue so a non-existent VM yields empty output
	// (treated as "not running") rather than an error.
	out, err := runPowerShell("(Get-VM -Name " + vmName + " -ErrorAction SilentlyContinue).State")
	if err != nil {
		return false, fmt.Errorf("not able to get VM state: %w", err)
	}
	return strings.EqualFold(strings.TrimSpace(out), "Running"), nil
}

func (h *hypervHypervisor) Pause() error {
	if _, err := runPowerShell("Suspend-VM -Name " + vmName + " -ErrorAction Stop"); err != nil {
		return fmt.Errorf("not able to pause VM: %w", err)
	}
	return nil
}

func (h *hypervHypervisor) Resume() error {
	if _, err := runPowerShell("Resume-VM -Name " + vmName + " -ErrorAction Stop"); err != nil {
		return fmt.Errorf("not able to resume VM: %w", err)
	}
	return nil
}

func (h *hypervHypervisor) TakeSnapshot(name string) error {
	cmd := fmt.Sprintf("Checkpoint-VM -Name %s -SnapshotName '%s' -ErrorAction Stop", vmName, name)
	if _, err := runPowerShell(cmd); err != nil {
		return fmt.Errorf("not able to take VM checkpoint: %w", err)
	}
	return nil
}

func (h *hypervHypervisor) DeleteSnapshot(name string) error {
	exists, err := h.snapshotExists(name)
	if err != nil {
		return fmt.Errorf("not able to check VM checkpoint: %w", err)
	}
	if !exists {
		return ErrSnapshotNotExist
	}
	cmd := fmt.Sprintf("Remove-VMSnapshot -VMName %s -Name '%s' -ErrorAction Stop", vmName, name)
	if _, err := runPowerShell(cmd); err != nil {
		return fmt.Errorf("not able to delete VM checkpoint: %w", err)
	}
	return nil
}

func (h *hypervHypervisor) RestoreSnapshot(name string) error {
	cmd := fmt.Sprintf("Restore-VMSnapshot -VMName %s -Name '%s' -Confirm:$false -ErrorAction Stop", vmName, name)
	if _, err := runPowerShell(cmd); err != nil {
		return fmt.Errorf("not able to restore VM checkpoint: %w", err)
	}
	return nil
}

// ResetNetworkLeases is a no-op on Hyper-V: there is no host-only network with
// persisted DHCP leases to clear (the virtual switch manages addressing).
func (h *hypervHypervisor) ResetNetworkLeases(hostOnlyCIDR string, verbose bool) error {
	return nil
}

// ApplyVB7Workaround is a no-op on Hyper-V (the VirtualBox 7 NAT workaround does
// not apply).
func (h *hypervHypervisor) ApplyVB7Workaround() error {
	return nil
}

// AddSwapDisk creates a dynamic VHDX and attaches it to the minikube VM.
//
// EXPERIMENTAL: attaching the disk to the correct controller depends on the VM
// generation, and the resulting in-VM device node is detected by the caller
// (see cmd swap helper) rather than assumed.
func (h *hypervHypervisor) AddSwapDisk(swapMB int16) error {
	swapPath := filepath.Join(utils.GetUserHome(), ".minikube", "machines", vmName, "swapdisk.vhdx")
	sizeBytes := int64(swapMB) * 1024 * 1024

	create := fmt.Sprintf("if (-not (Test-Path '%s')) { New-VHD -Path '%s' -SizeBytes %d -Dynamic | Out-Null }",
		swapPath, swapPath, sizeBytes)
	if _, err := runPowerShell(create); err != nil {
		return fmt.Errorf("cannot create swap VHD: %w", err)
	}

	// Gen 2 VMs use a SCSI controller (the default for Add-VMHardDiskDrive);
	// Gen 1 VMs only have IDE, so attach to a free IDE slot.
	attach := fmt.Sprintf("$gen = (Get-VM -Name %s).Generation; "+
		"if ($gen -eq 2) { Add-VMHardDiskDrive -VMName %s -Path '%s' -ErrorAction Stop } "+
		"else { Add-VMHardDiskDrive -VMName %s -ControllerType IDE -ControllerNumber 0 -ControllerLocation 1 -Path '%s' -ErrorAction Stop }",
		vmName, vmName, swapPath, vmName, swapPath)
	if _, err := runPowerShell(attach); err != nil {
		return fmt.Errorf("cannot attach swap VHD: %w", err)
	}
	return nil
}

func (h *hypervHypervisor) Validate(hypervVirtualSwitch string) error {
	elevated, err := isElevated()
	if err != nil {
		return err
	}
	if !elevated {
		return fmt.Errorf("the hyperv driver requires administrator privileges; please run gokube from an elevated (\"Run as administrator\") shell")
	}
	if _, err := runPowerShell("Get-Command Get-VM -ErrorAction Stop | Out-Null"); err != nil {
		return fmt.Errorf("Hyper-V does not appear to be enabled (Get-VM unavailable); enable it with \"Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All\" and reboot: %w", err)
	}
	if len(hypervVirtualSwitch) > 0 {
		cmd := fmt.Sprintf("Get-VMSwitch -Name '%s' -ErrorAction Stop | Out-Null", hypervVirtualSwitch)
		if _, err := runPowerShell(cmd); err != nil {
			return fmt.Errorf("hyperv virtual switch %q not found; create one or omit --hyperv-virtual-switch to use the Default Switch: %w", hypervVirtualSwitch, err)
		}
	}
	return nil
}

// snapshotExists reports whether a checkpoint with the given name exists for the
// minikube VM.
func (h *hypervHypervisor) snapshotExists(name string) (bool, error) {
	cmd := fmt.Sprintf("if (Get-VMSnapshot -VMName %s -Name '%s' -ErrorAction SilentlyContinue) { 'yes' } else { 'no' }", vmName, name)
	out, err := runPowerShell(cmd)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(out), "yes"), nil
}
