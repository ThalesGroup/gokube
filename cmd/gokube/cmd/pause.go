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
	"github.com/gemalto/gokube/pkg/virtualbox"
	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pauses minikube. This command pauses minikube VM",
	Long:  "Pauses minikube. This command pauses minikube VM",
	Run:   pauseRun,
}

func init() {
	RootCmd.AddCommand(pauseCmd)
}

func pauseRun(cmd *cobra.Command, args []string) {
	fmt.Println("Pausing minikube VM...")
	err := virtualbox.Pause()
	if err != nil {
		panic(err)
	}
}
