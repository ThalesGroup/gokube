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
	"github.com/cvila84/go-latest"
	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/spf13/cobra"
)

const (
	GOKUBE_VERSION = "1.13.0"
)

var gokubeVersion string

var (
	githubTag = &latest.GithubTag{
		Owner:      "ThalesGroup",
		Repository: "gokube",
	}
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version for gokube",
	Long:  `Shows version for gokube`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Fprintln(os.Stderr, "usage: gokube version")
			os.Exit(1)
		}
		fmt.Println("gokube version: v" + GOKUBE_VERSION)
		minikube.Version()
		docker.Version()
		kubectl.Version()
		helm.Version()
		helm.PluginsVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	res, _ := latest.Check(githubTag, GOKUBE_VERSION, 5*time.Second)
	if res == nil {
		fmt.Printf("WARNING: Cannot find gokube latest release, please check your connection\n")
	}
	if res != nil {
		if res.Outdated {
			fmt.Printf("WARNING: This version of gokube is outdated, please download the newest one on https://github.com/ThalesGroup/gokube/releases/tag/v%s\n", res.Current)
		} else if res.New {
			fmt.Printf("WARNING: This version of gokube has not yet been published, use it at your own risk !\n")
		}
	}
	gokube.ReadConfig()
	gokubeVersion = viper.GetString("gokube-version")
	if len(gokubeVersion) == 0 {
		gokubeVersion = "0.0.0"
	}
}
