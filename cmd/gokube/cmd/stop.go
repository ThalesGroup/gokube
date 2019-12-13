/*
(c) Copyright 2018, Gemalto. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
	"time"
)

var forceStop bool

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops minikube. This command stops minikube",
	Long:  "Stops minikube. This command stops minikube",
	Run:   stopRun,
}

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&forceStop, "force", "f", false, "Force VM stop (don't display warning message before stopping)")
}

func stopRun(cmd *cobra.Command, args []string) {
	if !forceStop {
		fmt.Println("WARNING: You should not stop a VM with a lot of running pods as the restart will be unstable")
		fmt.Print("Press <CTRL+C> within the next 10s it you need to perform some clean or press <ENTER> now to continue...")
		enter := make(chan bool, 1)
		go gokube.WaitEnter(enter)
		select {
		case <-enter:
		case <-time.After(10 * time.Second):
			fmt.Println()
		}
	}
	fmt.Println("Stopping minikube VM...")
	minikube.Stop()
}
