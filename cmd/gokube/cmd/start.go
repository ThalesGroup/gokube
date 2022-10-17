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
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:          "start",
	Short:        "Starts gokube. This command starts minikube",
	Long:         "Starts gokube. This command starts minikube",
	RunE:         startRun,
	SilenceUsage: true,
}

func init() {
	loadURLVersionsFromEnv()
	startCmd.Flags().BoolVarP(&askForUpgrade, "upgrade", "u", false, "Upgrade gokube (download and setup docker, minikube, kubectl and helm)")
	rootCmd.AddCommand(startCmd)
}

func start() error {
	err := gokube.ReadConfig(verbose)
	if err != nil {
		return fmt.Errorf("cannot read gokube configuration file: %w", err)
	}
	kubernetesVersionForStart := viper.GetString("kubernetes-version")
	if len(kubernetesVersionForStart) == 0 {
		kubernetesVersionForStart = utils.GetValueFromEnv("KUBERNETES_VERSION", DEFAULT_KUBERNETES_VERSION)
	}
	fmt.Printf("Starting minikube VM with kubernetes %s...\n", kubernetesVersionForStart)
	err = minikube.Restart(kubernetesVersionForStart)
	if err != nil {
		return fmt.Errorf("cannot restart minikube VM: %w", err)
	}
	return nil
}

func startRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}

	checkLatestVersion()

	if askForUpgrade {
		fmt.Println("Upgrading gokube dependencies...")
		err := upgradeDependencies()
		if err != nil {
			return err
		}
		fmt.Println("Upgrading helm plugins...")
		err = upgradeHelmPlugins()
		if err != nil {
			return err
		}
	}
	return start()
}
