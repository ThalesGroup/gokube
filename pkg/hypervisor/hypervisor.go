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

// Package hypervisor abstracts the host-side VM operations that minikube does
// not expose, so that gokube can drive either VirtualBox (via VBoxManage) or
// Hyper-V (via PowerShell) behind a single interface.
//
// Dependency rule: this package may import pkg/virtualbox but must NOT import
// pkg/minikube. The driver name and the Hyper-V virtual switch flow into
// minikube.Start as plain parameters from the cmd layer.
package hypervisor

import "errors"

const (
	// DriverVirtualBox is the default minikube driver (Oracle VirtualBox).
	DriverVirtualBox = "virtualbox"
	// DriverHyperV is the Microsoft Hyper-V minikube driver.
	DriverHyperV = "hyperv"
)

var (
	// ErrSnapshotNotExist is returned by DeleteSnapshot when the named snapshot
	// (VirtualBox) / checkpoint (Hyper-V) does not exist. It is a driver-neutral
	// sentinel so callers do not have to know which hypervisor produced it.
	ErrSnapshotNotExist = errors.New("snapshot does not exist")
	// ErrUnsupportedDriver is returned by New for an unknown driver name.
	ErrUnsupportedDriver = errors.New("unsupported driver, expected \"virtualbox\" or \"hyperv\"")
)

// Hypervisor is the set of host-side VM operations gokube needs on top of
// minikube. Implementations operate on the "minikube" VM by name.
type Hypervisor interface {
	// IsRunning reports whether the minikube VM is currently running.
	IsRunning() (bool, error)
	// Pause suspends the running VM.
	Pause() error
	// Resume wakes a paused/suspended VM.
	Resume() error
	// TakeSnapshot creates a snapshot/checkpoint with the given name.
	TakeSnapshot(name string) error
	// DeleteSnapshot removes the named snapshot/checkpoint. It returns
	// ErrSnapshotNotExist if no such snapshot exists.
	DeleteSnapshot(name string) error
	// RestoreSnapshot restores the VM to the named snapshot/checkpoint.
	RestoreSnapshot(name string) error
	// ResetNetworkLeases clears stale host-only DHCP leases so the VM can get
	// its expected IP. No-op for drivers without host-only networking (Hyper-V).
	ResetNetworkLeases(hostOnlyCIDR string, verbose bool) error
	// ApplyVB7Workaround applies the VirtualBox 7 NAT localhost-reachable
	// workaround. No-op for drivers other than VirtualBox.
	ApplyVB7Workaround() error
	// AddSwapDisk creates and attaches a host-side swap disk of swapMB megabytes
	// to the minikube VM. In-VM formatting/activation is handled by the caller.
	AddSwapDisk(swapMB int16) error
	// Validate checks that the host is ready for this driver (e.g. Hyper-V
	// enabled, process elevated, virtual switch present). hypervVirtualSwitch
	// may be empty to let minikube auto-select a switch.
	Validate(hypervVirtualSwitch string) error
}

// New returns the Hypervisor implementation for the given driver name.
func New(driver string) (Hypervisor, error) {
	switch driver {
	case DriverVirtualBox:
		return &vboxHypervisor{}, nil
	case DriverHyperV:
		return &hypervHypervisor{}, nil
	default:
		return nil, ErrUnsupportedDriver
	}
}
