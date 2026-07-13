//go:build !windows

package virtualbox

import "fmt"

func findVBoxInstallDirInRegistry() (string, error) {
	return "", fmt.Errorf("VirtualBox registry lookup is only available on Windows")
}
