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

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:          "stop",
	Short:        "Stops gokube. This command stops minikube",
	Long:         "Stops gokube. This command stops minikube",
	RunE:         stopRun,
	SilenceUsage: true,
}

func init() {
	defaultGokubeQuiet := false
	if len(getValueFromEnv("GOKUBE_QUIET", "")) > 0 {
		defaultGokubeQuiet = true
	}
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&quiet, "quiet", "q", defaultGokubeQuiet, "Don't display warning message before stopping")
}

func confirmStopCommandExecution() {
	fmt.Println("WARNING: You should not stop a VM with a lot of running pods as the restart will be unstable")
	fmt.Print("Press <CTRL+C> within the next 10s it you need to perform some clean or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go gokube.WaitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func stopRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}
	if !quiet {
		confirmStopCommandExecution()
	}
	fmt.Println("Stopping minikube VM...")
	return minikube.Stop()
}
