package virtualbox

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func findVBoxInstallDirInRegistry() (string, error) {
	registryKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Oracle\VirtualBox`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("can't find VirtualBox registry entries, is VirtualBox really installed properly? %w", err)
	}
	defer registryKey.Close()

	installDir, _, err := registryKey.GetStringValue("InstallDir")
	if err != nil {
		return "", fmt.Errorf("can't find InstallDir registry key within VirtualBox registries entries, is VirtualBox really installed properly? %w", err)
	}

	return installDir, nil
}
