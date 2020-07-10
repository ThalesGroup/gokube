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
	"fmt"
	"github.com/coreos/go-semver/semver"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/gokube"
	"github.com/gemalto/gokube/pkg/utils"
)

// Start ...
func Start(memory int16, cpus int16, diskSize string, httpProxy string, httpsProxy string, noProxy string, insecureRegistry string, kubernetesVersion string, cache bool, dnsProxy bool, hostDNSResolver bool) {
	var args = []string{"start", "--kubernetes-version", kubernetesVersion, "--insecure-registry", insecureRegistry, "--memory", strconv.FormatInt(int64(memory), 10), "--cpus", strconv.FormatInt(int64(cpus), 10), "--disk-size", diskSize, "--network-plugin=cni", "--cni=bridge"}
	//patchStartArgs(args, kubernetesVersion)
	if len(httpProxy) > 0 {
		args = append(args, "--docker-env=http_proxy="+httpProxy)
	}
	if len(httpsProxy) > 0 {
		args = append(args, "--docker-env=https_proxy="+httpsProxy)
	}
	if len(noProxy) > 0 {
		args = append(args, "--docker-env=no_proxy="+noProxy)
	}
	if !cache {
		args = append(args, "--cache-images=false")
	}
	if dnsProxy {
		args = append(args, "--dns-proxy")
	}
	if !hostDNSResolver {
		args = append(args, "--host-dns-resolver=false")
	}
	cmd := exec.Command("minikube", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Restart ...
func Restart(kubernetesVersion string) {
	var args = []string{"start", "--kubernetes-version", kubernetesVersion}
	//patchStartArgs(args, kubernetesVersion)
	cmd := exec.Command("minikube", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Stop ...
func Stop() {
	cmd := exec.Command("minikube", "stop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Delete ...
func Delete() {
	cmd := exec.Command("minikube", "delete")
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
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
		fmt.Print(err)
		os.Exit(1)
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

// IP...
func Ip() string {
	out, err := exec.Command("minikube", "ip").Output()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	return strings.TrimRight(string(out), "\r\n")
}

// DownloadExecutable ...
func DownloadExecutable(dst string, minikubeURI string, minikubeVersion string) {
	if _, err := os.Stat(dst + "/minikube.exe"); os.IsNotExist(err) {
		download.FromUrl("minikube "+minikubeVersion, minikubeURI, minikubeVersion)
		utils.MoveFile(gokube.GetTempDir()+"/minikube-windows-amd64.exe", dst+"/minikube.exe")
		utils.RemoveDir(gokube.GetTempDir())
	}
}

// DeleteExecutable ...
func DeleteExecutable() {
	utils.RemoveFile(gokube.GetBinDir() + "/minikube.exe")
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() {
	utils.CleanDir(utils.GetUserHome() + "/.minikube")
}

func patchStartArgs(args []string, kubernetesVersion string) {
	if semver.New(kubernetesVersion[1:]).Compare(*semver.New("1.16.0")) >= 0 && semver.New(kubernetesVersion[1:]).Compare(*semver.New("1.18.0")) < 0 {
		args = append(args, "--extra-config=apiserver.runtime-config=apps/v1beta1=true,apps/v1beta2=true,extensions/v1beta1/daemonsets=true,extensions/v1beta1/deployments=true,extensions/v1beta1/replicasets=true,extensions/v1beta1/networkpolicies=true,extensions/v1beta1/podsecuritypolicies=true")
	}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
