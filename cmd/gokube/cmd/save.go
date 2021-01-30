package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
	"time"
)

var live bool

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
	if len(getValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	saveCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before snapshotting")
	saveCmd.Flags().BoolVarP(&live, "live", "l", false, "Don't stop VM before taking snapshot")
	rootCmd.AddCommand(saveCmd)
}

func confirmSnapshotCommandExecution() {
	fmt.Println("WARNING: You should not snapshot a running VM as the process can be long and take more space on disk")
	fmt.Print("Press <CTRL+C> within the next 10s it you want to stop VM first or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go gokube.WaitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func saveRun(cmd *cobra.Command, args []string) error {
	var running bool
	var err error
	if len(args) > 0 {
		return cmd.Usage()
	}
	if live && !quiet {
		confirmSnapshotCommandExecution()
	} else if !live {
		running, err = virtualbox.IsRunning()
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
	}
	fmt.Println("Taking snapshot of minikube VM...")
	_ = virtualbox.DeleteSnapshot("gokube")
	err = virtualbox.TakeSnapshot("gokube")
	if err != nil {
		return err
	}
	if running {
		return start()
	} else {
		return nil
	}
}
