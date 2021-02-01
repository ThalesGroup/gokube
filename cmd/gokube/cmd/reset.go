package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
)

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
	if len(getValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	resetCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before resetting")
	rootCmd.AddCommand(resetCmd)
}

func resetRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}
	running, err := virtualbox.IsRunning()
	if err != nil {
		return err
	}
	if running {
		fmt.Println("Stopping minikube VM...")
		err = minikube.Stop()
		if err != nil {
			return err
		}
	}
	fmt.Println("Resetting minikube VM from snapshot...")
	err = virtualbox.RestoreSnapshot("gokube")
	if err != nil {
		return err
	}
	if running {
		return start()
	} else {
		return nil
	}
}
