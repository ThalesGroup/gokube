package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/gemalto/gokube/pkg/virtualbox"
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

	running, err := virtualbox.IsRunning()
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
	err = virtualbox.RestoreSnapshot(snapshotName)
	if err != nil {
		return fmt.Errorf("cannot restore minikube VM snapshot %s: %w", snapshotName, err)
	}
	if clean {
		err = virtualbox.DeleteSnapshot(snapshotName)
		if err != nil && err != virtualbox.ErrSnapshotNotExist {
			return fmt.Errorf("cannot delete minikube VM snapshot %s: %w", snapshotName, err)
		}
	}
	fmt.Printf("Minikube VM has successfully been reset from snapshot '%s'\n", snapshotName)
	if running {
		return start()
	} else {
		return nil
	}
}
