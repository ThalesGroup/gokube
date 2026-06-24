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
	"os/exec"
	"strings"
	"time"
)

// swapDetectRetries / swapDetectDelay bound how long we wait for a freshly
// attached disk to appear inside the running VM (hot-plug detection is
// asynchronous).
const (
	swapDetectRetries = 6
	swapDetectDelay   = 1 * time.Second
)

// runMinikubeSSH runs a command inside the minikube VM and returns its trimmed
// stdout.
func runMinikubeSSH(command string) (string, error) {
	out, err := exec.Command("minikube", "ssh", command).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}

// listMinikubeDisks returns the set of whole-disk device names (e.g. "sda",
// "sdb") currently visible inside the minikube VM.
func listMinikubeDisks() (map[string]bool, error) {
	out, err := runMinikubeSSH("lsblk -dn -o NAME")
	if err != nil {
		return nil, err
	}
	disks := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		name := strings.TrimSpace(line)
		if len(name) > 0 {
			disks[name] = true
		}
	}
	return disks, nil
}

// detectNewSwapDevice returns the /dev path of the disk that appeared after the
// swap disk was attached, by diffing against the set of disks present before.
// It retries a few times because hot-plug detection inside the VM is
// asynchronous.
func detectNewSwapDevice(before map[string]bool) (string, error) {
	for attempt := 0; attempt < swapDetectRetries; attempt++ {
		after, err := listMinikubeDisks()
		if err != nil {
			return "", err
		}
		var added []string
		for name := range after {
			if !before[name] {
				added = append(added, name)
			}
		}
		switch {
		case len(added) == 1:
			return "/dev/" + added[0], nil
		case len(added) > 1:
			return "", fmt.Errorf("found %d new disks (%s), cannot determine which is the swap disk", len(added), strings.Join(added, ", "))
		}
		time.Sleep(swapDetectDelay)
	}
	return "", fmt.Errorf("no new disk appeared in the minikube VM after attaching swap")
}

// formatAndEnableSwap formats the given device as swap, activates it, and adds a
// persistent /etc/fstab entry keyed by UUID (robust against device-node
// reordering across reboots). The fstab entry lets `gokube start` simply run
// `swapon --all`.
func formatAndEnableSwap(device string) error {
	if _, err := runMinikubeSSH("sudo mkswap " + device); err != nil {
		return fmt.Errorf("cannot mkswap %s: %w", device, err)
	}

	// Resolve the swap UUID so the fstab entry is stable. Fall back to the
	// device node if the UUID cannot be read.
	uuid, err := runMinikubeSSH("sudo blkid -s UUID -o value " + device)
	if err != nil || len(strings.TrimSpace(uuid)) == 0 {
		uuid = ""
	}
	uuid = strings.TrimSpace(uuid)

	var spec, fstabLine string
	if len(uuid) > 0 {
		spec = "UUID=" + uuid
		fstabLine = "UUID=" + uuid + " none swap defaults 0 0"
	} else {
		spec = device
		fstabLine = device + " none swap defaults 0 0"
	}

	if _, err := runMinikubeSSH("sudo swapon " + spec); err != nil {
		return fmt.Errorf("cannot swapon %s: %w", spec, err)
	}
	if _, err := runMinikubeSSH(fmt.Sprintf("echo '%s' | sudo tee -a /etc/fstab", fstabLine)); err != nil {
		return fmt.Errorf("cannot persist swap in /etc/fstab: %w", err)
	}
	return nil
}
