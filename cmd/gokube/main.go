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

package main

import (
	"fmt"
	"github.com/gemalto/gokube/cmd/gokube/cmd"
	"github.com/tcnksm/go-latest"
)

var (
	githubTag = &latest.GithubTag{
		Owner:      "ThalesGroup",
		Repository: "gokube",
	}
)

func main() {
	res, _ := latest.Check(githubTag, cmd.GOKUBE_VERSION)
	if res == nil {
		fmt.Printf("WARNING: Cannot look for gokube upgrades, please check your connection\n")
	}
	if res != nil && res.Outdated {
		fmt.Printf("WARNING: This version of gokube is outdated, please download the newest one on https://github.com/ThalesGroup/gokube/releases/tag/v%s\n", res.Current)
	}
	cmd.Execute()
}
