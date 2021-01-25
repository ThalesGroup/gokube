package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
)

// resetCmd represents the pause command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset minikube from previously taken snapshot.",
	Long:  "Reset minikube from previously taken snapshot.",
	Run:   resetRun,
}

func init() {
	RootCmd.AddCommand(resetCmd)
}

func resetRun(cmd *cobra.Command, args []string) {
	stopRun(cmd, args)
	fmt.Println("Resetting minikube VM from snapshot...")
	err := virtualbox.RestoreSnapshot("gokube")
	if err != nil {
		panic(err)
	}
	startRun(cmd, args)
}
