package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
)

// snapshotCmd represents the pause command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Snapshot minikube VM to set the reference for reset command.",
	Long:  "Snapshot minikube VM to set the reference for reset command.",
	Run:   snapshotRun,
}

func init() {
	RootCmd.AddCommand(snapshotCmd)
}

func snapshotRun(cmd *cobra.Command, args []string) {
	fmt.Println("Taking snapshot of minikube VM...")
	err := virtualbox.TakeSnapshot("gokube")
	if err != nil {
		panic(err)
	}
}
