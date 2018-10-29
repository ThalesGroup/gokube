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

package minikube

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

// Start ...
func Start(memory int16, cpus int16, diskSize string, tproxy bool, httpProxy string, httpsProxy string, noProxy string, insecureRegistry string, kubernetesVersion string, cache bool) {
	var args = []string{"start", "--kubernetes-version", kubernetesVersion, "--insecure-registry", insecureRegistry, "--memory", strconv.FormatInt(int64(memory), 10), "--cpus", strconv.FormatInt(int64(cpus), 10), "--disk-size", diskSize, "--network-plugin=cni", "--extra-config=kubelet.network-plugin=cni"}
	if !tproxy {
		args = append(args, "--docker-env", "HTTP_PROXY="+httpProxy, "--docker-env", "HTTPS_PROXY="+httpsProxy, "--docker-env", "NO_PROXY="+noProxy)
	}
	if cache {
		args = append(args, "--cache-images")
	}
	cmd := exec.Command("minikube", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Cache ...
func Cache(image string) {
	cmd := exec.Command("minikube", "cache", "add", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// AddonsEnable ...
func AddonsEnable(addon string) {
	cmd := exec.Command("minikube", "addons", "enable", addon)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// QuickStart ...
func QuickStart() {
	cmd := exec.Command("minikube", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Stop ...
func Stop() {
	cmd := exec.Command("minikube", "stop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Status ...
func Status() {
	cmd := exec.Command("minikube", "status")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Delete ...
func Delete() {
	cmd := exec.Command("minikube", "delete")
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	cmd.Run()
}

// CopyCerts ...
func CopyCerts() {
	path := toLinuxPath("cp -r /home/docker/* etc/docker/certs.d")
	cmd := exec.Command("minikube", "ssh", "sudo", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// ConfigSet ...
func ConfigSet(key string, value string) {
	cmd := exec.Command("minikube", "config", "set", key, value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Version ...
func Version() {
	cmd := exec.Command("minikube", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// DockerEnv ...
func DockerEnv() []utils.EnvVar {
	out, err := exec.Command("minikube", "docker-env").Output()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var envVar []utils.EnvVar
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "export ") {
			tokens := strings.Split(strings.TrimPrefix(line, "export "), "=")
			envVar = append(envVar, utils.EnvVar{
				Name:  tokens[0],
				Value: trimQuotes(tokens[1]),
			})
		}
		if strings.HasPrefix(line, "SET ") {
			tokens := strings.Split(strings.TrimPrefix(line, "SET "), "=")
			envVar = append(envVar, utils.EnvVar{
				Name:  tokens[0],
				Value: trimQuotes(tokens[1]),
			})
		}
	}
	return envVar
}

// Download ...
func Download(dst string, minikubeURI string, minikubeVersion string) {
	if _, err := os.Stat(dst + "/minikube.exe"); os.IsNotExist(err) {
		download.DownloadFromUrl("minikube "+minikubeVersion, minikubeURI, minikubeVersion)
		utils.MoveFile(gokube.GetTempDir()+"/minikube-windows-amd64.exe", dst+"/minikube.exe")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// Purge ...
func Purge() {
	utils.RemoveFile(gokube.GetBinDir() + "/minikube.exe")
	utils.CleanDir(utils.GetUserHome() + "/.minikube")
}

func toLinuxPath(path string) string {
	path = strings.Replace(path, "\\", "/", -1)
	return strings.Replace(path, "C:", "/c", -1)
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
