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
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// GetBinDir ...
func GetBinDir() string {
	path, err := exec.LookPath("gokube")
	if err != nil {
		panic(err)
	}
	if path == "gokube.exe" {
		path = whereAmI()
	} else {
		path = strings.TrimSuffix(path, "\\gokube.exe")
	}
	return path
}

// GetTempDir ...
func GetTempDir() string {
	return GetBinDir() + "/tmp"
}

// ReadConfig ...
func ReadConfig() {
	usr, err := user.Current()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	configPath := usr.HomeDir + "/.gokube"
	configFile := "config"
	configFilePath := configPath + "/config.yaml"
	if _, existDirErr := os.Stat(configPath); os.IsNotExist(existDirErr) {
		createDirErr := os.Mkdir(configPath, os.ModePerm)
		if createDirErr != nil {
			fmt.Print(createDirErr)
			os.Exit(1)
		}
	}
	if _, existFileErr := os.Stat(configFilePath); os.IsNotExist(existFileErr) {
		_, createFileErr := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
		if createFileErr != nil {
			fmt.Print(createFileErr)
			os.Exit(1)
		}
	}
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		fmt.Println(readConfigErr)
	}
}

// WriteConfig ...
func WriteConfig(gokubeVersion string, kubernetesVersion string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	configPath := usr.HomeDir + "/.gokube"
	configFile := "config"
	configFilePath := configPath + "/config.yaml"
	if _, existDirErr := os.Stat(configPath); os.IsNotExist(existDirErr) {
		createDirErr := os.Mkdir(configPath, os.ModePerm)
		if createDirErr != nil {
			fmt.Print(createDirErr)
			os.Exit(1)
		}
	}
	if _, existFileErr := os.Stat(configFilePath); os.IsNotExist(existFileErr) {
		_, createFileErr := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
		if createFileErr != nil {
			fmt.Print(createFileErr)
			os.Exit(1)
		}
	}
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.Set("gokube-version", gokubeVersion)
	viper.Set("kubernetes-version", kubernetesVersion)
	err = viper.WriteConfig()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

// WaitEnter...
func WaitEnter(enter chan<- bool) {
	_, _, _ = bufio.NewReader(os.Stdin).ReadLine()
	enter <- true
}

// WhereAmI returns a string containing the file name, function name
// and the line number of a specified entry on the call stack
func whereAmI() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	return dir
}
