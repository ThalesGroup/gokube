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

package gokube

import (
	"bufio"
	"fmt"
	"github.com/gemalto/gokube/pkg/docker"
	"github.com/gemalto/gokube/pkg/helm"
	"github.com/gemalto/gokube/pkg/helmimage"
	"github.com/gemalto/gokube/pkg/helmpush"
	"github.com/gemalto/gokube/pkg/helmspray"
	"github.com/gemalto/gokube/pkg/kubectl"
	"github.com/gemalto/gokube/pkg/minikube"
	"github.com/gemalto/gokube/pkg/stern"
	"github.com/gemalto/gokube/pkg/k9s"
	"github.com/gemalto/gokube/pkg/utils"
	"github.com/spf13/viper"
	"os"
	"time"
)

type HelmPlugins struct {
	SprayURL     string
	SprayVersion string
	ImageURL     string
	ImageVersion string
	PushURL      string
	PushVersion  string
}

type Dependencies struct {
	MinikubeURL     string
	MinikubeVersion string
	HelmURL         string
	HelmVersion     string
	DockerURL       string
	DockerVersion   string
	KubectlURL      string
	KubectlVersion  string
	SternURL        string
	SternVersion    string
	K9sURL          string
	K9sVersion      string
}

// ReadConfig ...
func ReadConfig(verbose bool) error {
	configPath := utils.GetUserHome() + string(os.PathSeparator) + ".gokube"
	if verbose {
		fmt.Printf("Reading %s...\n", configPath)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	configFilePath := configPath + string(os.PathSeparator) + "config.yaml"
	if verbose {
		fmt.Printf("Reading %s...\n", configFilePath)
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		_, err = os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if verbose {
		fmt.Printf("Read settings: %+v\n", viper.AllSettings())
	}
	if err != nil {
		return err
	}
	return nil
}

// WriteConfig ...
func WriteConfig(gokubeVersion string, kubernetesVersion string, containerRuntime string) error {
	configPath := utils.GetUserHome() + string(os.PathSeparator) + ".gokube"
	configFile := "config"
	configFilePath := configPath + string(os.PathSeparator) + "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		_, err = os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	}
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.Set("gokube-version", gokubeVersion)
	viper.Set("kubernetes-version", kubernetesVersion)
	viper.Set("container-runtime", containerRuntime)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func UpgradeHelmPlugins(plugins *HelmPlugins) error {
	// TODO rely on helm plugin install
	_ = helmspray.DeletePlugin()
	err := helmspray.InstallPlugin(plugins.SprayURL, plugins.SprayVersion)
	if err != nil {
		return fmt.Errorf("cannot install helm-spray plugin: %w", err)
	}
	_ = helmimage.DeletePlugin()
	err = helmimage.InstallPlugin(plugins.ImageURL, plugins.ImageVersion)
	if err != nil {
		return fmt.Errorf("cannot install helm-image plugin: %w", err)
	}
	_ = helmpush.DeletePlugin()
	err = helmpush.InstallPlugin(plugins.PushURL, plugins.PushVersion)
	if err != nil {
		return fmt.Errorf("cannot install helm-push plugin: %w", err)
	}
	return nil
}

func UpgradeDependencies(dependencies *Dependencies) error {
	_ = minikube.DeleteExecutable()
	err := minikube.DownloadExecutable(dependencies.MinikubeURL, dependencies.MinikubeVersion)
	if err != nil {
		return fmt.Errorf("cannot download or upgrade minikube: %w", err)
	}
	_ = helm.DeleteExecutable()
	err = helm.DownloadExecutable(dependencies.HelmURL, dependencies.HelmVersion)
	if err != nil {
		return fmt.Errorf("cannot download or upgrade helm: %w", err)
	}
	_ = docker.DeleteExecutable()
	err = docker.DownloadExecutable(dependencies.DockerURL, dependencies.DockerVersion)
	if err != nil {
		return fmt.Errorf("cannot download or upgrade docker: %w", err)
	}
	_ = kubectl.DeleteExecutable()
	err = kubectl.DownloadExecutable(dependencies.KubectlURL, dependencies.KubectlVersion)
	if err != nil {
		return fmt.Errorf("cannot download or upgrade kubectl: %w", err)
	}
	_ = stern.DeleteExecutable()
	err = stern.DownloadExecutable(dependencies.SternURL, dependencies.SternVersion)
	if err != nil {
		return fmt.Errorf("cannot download or upgrade stern: %w", err)
	}
	_ = k9s.DeleteExecutable()
	err = k9s.DownloadExecutable(dependencies.K9sURL, dependencies.K9sVersion)
	if err != nil {
		return fmt.Errorf("cannot download or upgrade k9s: %w", err)
	}
	return nil
}

func ConfirmInitCommandExecution() {
	fmt.Println("Warning: your Virtualbox GUI shall not be open and no other VM shall be currently running")
	fmt.Print("Press <CTRL+C> within the next 10s it you need to check this or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go waitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func ConfirmSnapshotCommandExecution() {
	fmt.Println("Warning: you should not snapshot a running VM as the process can be long and take more space on disk")
	fmt.Print("Press <CTRL+C> within the next 10s it you want to stop VM first or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go waitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func ConfirmStopCommandExecution() {
	fmt.Println("Warning: you should not stop a VM with a lot of running pods as the restart will be unstable")
	fmt.Print("Press <CTRL+C> within the next 10s it you need to perform some clean or press <ENTER> now to continue...")
	enter := make(chan bool, 1)
	go waitEnter(enter)
	select {
	case <-enter:
	case <-time.After(10 * time.Second):
		fmt.Println()
	}
	time.Sleep(200 * time.Millisecond)
}

func waitEnter(enter chan<- bool) {
	_, _, _ = bufio.NewReader(os.Stdin).ReadLine()
	enter <- true
}
