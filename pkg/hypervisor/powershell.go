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
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// powerShellCmd is the PowerShell executable used to drive Hyper-V. Windows
// PowerShell ships with the Hyper-V module and is always on PATH on hosts where
// Hyper-V is available.
const powerShellCmd = "powershell"

// runPowerShell executes a single PowerShell command line and returns its
// trimmed stdout. On failure it returns an error that includes stderr so the
// underlying cmdlet message is preserved.
func runPowerShell(command string) (string, error) {
	cmd := exec.Command(powerShellCmd, "-NoProfile", "-NonInteractive", "-Command", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := strings.TrimRight(stdout.String(), "\r\n")
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if len(errMsg) == 0 {
			errMsg = err.Error()
		}
		return out, fmt.Errorf("powershell %q failed: %s", command, errMsg)
	}
	return out, nil
}
