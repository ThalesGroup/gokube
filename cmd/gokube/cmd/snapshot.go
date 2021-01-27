package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
	"time"
)

// snapshotCmd represents the pause command
var snapshotCmd = &cobra.Command{
	Use:          "snapshot",
	Short:        "Creates a gokube reference. This command takes a snapshot of the minikube VM (which will be the reference for reset command)",
	Long:         "Creates a gokube reference. This command takes a snapshot of the minikube VM (which will be the reference for reset command)",
	RunE:         snapshotRun,
	SilenceUsage: true,
}

func init() {
	defaultGokubeQuiet := false
	if len(getValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	RootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before snapshotting")
}

func confirmSnapshotCommandExecution() {
	fmt.Println("WARNING: You should not snapshot a running VM as the process can be long and take more space on disk")
	fmt.Print("Press <CTRL+C> within the next 10s it you need to stop VM or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go gokube.WaitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func snapshotRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}
	if !quiet {
		confirmSnapshotCommandExecution()
	}
	fmt.Println("Taking snapshot of minikube VM...")
	return virtualbox.TakeSnapshot("gokube")
}
