package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
)

var live bool
var savedSnapshotName string

// saveCmd represents the pause command
var saveCmd = &cobra.Command{
	Use:          "save",
	Short:        "Creates a gokube reference. This command takes a snapshot of the minikube VM (which will be the reference for reset command)",
	Long:         "Creates a gokube reference. This command takes a snapshot of the minikube VM (which will be the reference for reset command)",
	RunE:         saveRun,
	SilenceUsage: true,
}

func init() {
	defaultGokubeQuiet := false
	if len(utils.GetValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	saveCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before snapshotting")
	saveCmd.Flags().BoolVarP(&live, "live", "l", false, "Don't stop VM before taking snapshot")
	saveCmd.Flags().StringVarP(&savedSnapshotName, "name", "n", "gokube", "The snapshot name")
	rootCmd.AddCommand(saveCmd)
}

func saveRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}

	checkLatestVersion()

	running := false
	if live && !quiet {
		gokube.ConfirmSnapshotCommandExecution()
	} else if !live {
		var err error
		running, err = virtualbox.IsRunning()
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
	}
	fmt.Printf("Taking snapshot '%s' of minikube VM...\n", savedSnapshotName)
	err := virtualbox.DeleteSnapshot(savedSnapshotName)
	if err != nil {
		return fmt.Errorf("cannot delete minikube VM snapshot %s: %w", savedSnapshotName, err)
	}
	err = virtualbox.TakeSnapshot(savedSnapshotName)
	if err != nil {
		return fmt.Errorf("cannot take minikube VM snapshot %s: %w", savedSnapshotName, err)
	}
	if running {
		return start()
	} else {
		return nil
	}
}
