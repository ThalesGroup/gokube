package cmd

import (
	"errors"
	"fmt"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/hypervisor"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/spf13/cobra"
)

var clean bool

// resetCmd represents the pause command
var resetCmd = &cobra.Command{
	Use:          "reset",
	Short:        "Resets gokube. This command restores minikube VM from previously taken snapshot",
	Long:         "Resets gokube. This command restores minikube VM from previously taken snapshot",
	RunE:         resetRun,
	SilenceUsage: true,
}

func init() {
	defaultGokubeQuiet := false
	if len(utils.GetValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	resetCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before resetting")
	resetCmd.Flags().StringVarP(&snapshotName, "name", "n", "gokube", "The snapshot name")
	resetCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Clean snapshot after reset")
	rootCmd.AddCommand(resetCmd)
}

func resetRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}

	checkLatestVersion()

	if err := gokube.ReadConfig(verbose); err != nil {
		return fmt.Errorf("cannot read gokube configuration file: %w", err)
	}
	hv, err := hypervisor.New(resolveDriver())
	if err != nil {
		return fmt.Errorf("invalid minikube driver: %w", err)
	}

	running, err := hv.IsRunning()
	if err != nil {
		return fmt.Errorf("cannot check if minikube VM is running: %w", err)
	}
	if running {
		fmt.Println("Stopping minikube VM...")
		err = minikube.Stop()
		if err != nil {
			return fmt.Errorf("cannot stop minikube VM: %w", err)
		}
	}
	fmt.Printf("Resetting minikube VM from snapshot '%s'...\n", snapshotName)
	err = hv.RestoreSnapshot(snapshotName)
	if err != nil {
		return fmt.Errorf("cannot restore minikube VM snapshot %s: %w", snapshotName, err)
	}
	if clean {
		err = hv.DeleteSnapshot(snapshotName)
		if err != nil && !errors.Is(err, hypervisor.ErrSnapshotNotExist) {
			return fmt.Errorf("cannot delete minikube VM snapshot %s: %w", snapshotName, err)
		}
	}
	fmt.Printf("Minikube VM has successfully been reset from snapshot '%s'\n", snapshotName)
	if running {
		fmt.Println("VM was running before reset, restarting...")
	} else {
		fmt.Println("VM was stopped before reset, starting after restore...")
	}
	return start()
}
