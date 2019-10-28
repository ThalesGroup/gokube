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
	"github.com/spf13/viper"
	"os"
)

var currentKubernetesVersion string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts minikube. This command starts minikube",
	Long:  "Starts minikube. This command starts minikube",
	Run:   startRun,
}

func init() {
	gokube.ReadConfig()
	defaultKubernetesVersion := viper.GetString("kubernetes-version")
	if len(defaultKubernetesVersion) == 0 {
		defaultKubernetesVersion = os.Getenv("KUBERNETES_VERSION")
		if len(defaultKubernetesVersion) == 0 {
			defaultKubernetesVersion = DEFAULT_KUBERNETES_VERSION
		}
	}
	startCmd.Flags().StringVarP(&currentKubernetesVersion, "kubernetes-version", "", defaultKubernetesVersion, "The kubernetes version")
	RootCmd.AddCommand(startCmd)
}

func startRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Starting minikube VM with kubernetes %s...\n", currentKubernetesVersion)
	minikube.Restart(currentKubernetesVersion)
}
