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

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:          "resume",
	Short:        "Resumes gokube. This command resumes the minikube VM",
	Long:         "Resumes gokube. This command resumes the minikube VM",
	RunE:         resumeRun,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func resumeRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Usage()
	}

	checkLatestVersion()

	fmt.Println("Resuming minikube VM...")
	err := virtualbox.Resume()
	if err != nil {
		return fmt.Errorf("cannot resume minikube VM: %w", err)
	}
	return nil
}
