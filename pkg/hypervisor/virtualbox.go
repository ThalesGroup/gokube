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

import "github.com/gemalto/gokube/pkg/virtualbox"

// vboxHypervisor implements Hypervisor by delegating to the existing
// pkg/virtualbox package (VBoxManage). It preserves the original behavior so
// existing VirtualBox users are unaffected.
type vboxHypervisor struct{}

func (v *vboxHypervisor) IsRunning() (bool, error) {
	return virtualbox.IsRunning()
}

func (v *vboxHypervisor) Pause() error {
	return virtualbox.Pause()
}

func (v *vboxHypervisor) Resume() error {
	return virtualbox.Resume()
}

func (v *vboxHypervisor) TakeSnapshot(name string) error {
	return virtualbox.TakeSnapshot(name)
}

func (v *vboxHypervisor) DeleteSnapshot(name string) error {
	err := virtualbox.DeleteSnapshot(name)
	// Translate the package-specific sentinel into the driver-neutral one so
	// callers can compare against hypervisor.ErrSnapshotNotExist.
	if err == virtualbox.ErrSnapshotNotExist {
		return ErrSnapshotNotExist
	}
	return err
}

func (v *vboxHypervisor) RestoreSnapshot(name string) error {
	return virtualbox.RestoreSnapshot(name)
}

func (v *vboxHypervisor) ResetNetworkLeases(hostOnlyCIDR string, verbose bool) error {
	return virtualbox.ResetHostOnlyNetworkLeases(hostOnlyCIDR, verbose)
}

func (v *vboxHypervisor) ApplyVB7Workaround() error {
	return virtualbox.Update("--nat-localhostreachable1=on")
}

func (v *vboxHypervisor) AddSwapDisk(swapMB int16) error {
	return virtualbox.NewVBoxManager().AddSwapDisk(swapMB)
}

func (v *vboxHypervisor) Validate(hypervVirtualSwitch string) error {
	// VirtualBox has no pre-flight validation here; minikube performs its own.
	return nil
}
