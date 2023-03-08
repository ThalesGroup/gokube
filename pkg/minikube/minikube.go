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
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gemalto/gokube/pkg/download"
	"github.com/gemalto/gokube/pkg/utils"
)

const (
	DEFAULT_URL           = "https://storage.googleapis.com/minikube/releases/%s/minikube-windows-amd64.exe"
	LOCAL_EXECUTABLE_NAME = "minikube.exe"
)

// Start ...
func Start(memory int16, cpus int16, diskSize string, httpProxy string, httpsProxy string, noProxy string, insecureRegistry string, kubernetesVersion string, cache bool, dnsProxy bool, hostDNSResolver bool, dnsDomain string, containerRuntime string, force bool, verbose bool) error {
	var args = []string{"start", "--kubernetes-version", kubernetesVersion, "--insecure-registry", insecureRegistry, "--memory", strconv.FormatInt(int64(memory), 10), "--cpus", strconv.FormatInt(int64(cpus), 10), "--disk-size", diskSize, "--driver=virtualbox", "--host-only-cidr=192.168.99.1/24"}
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
	if len(dnsDomain) > 0 {
		args = append(args, "--dns-domain="+dnsDomain)
	}
	if len(containerRuntime) > 0 {
		args = append(args, "--container-runtime="+containerRuntime)
	}
	if force {
		args = append(args, "--force")
	}
	if verbose {
		args = append(args, "--alsologtostderr", "--v=1")
	}
	cmd := exec.Command("minikube", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Restart ...
func Restart(kubernetesVersion string, containerRuntime string, force bool, verbose bool) error {
	var args = []string{"start", "--kubernetes-version", kubernetesVersion}
	if len(containerRuntime) > 0 {
		args = append(args, "--container-runtime="+containerRuntime)
	}
	if force {
		args = append(args, "--force")
	}
	if verbose {
		args = append(args, "--alsologtostderr", "--v=1")
	}
	cmd := exec.Command("minikube", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Stop ...
func Stop() error {
	cmd := exec.Command("minikube", "stop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Delete ...
func Delete() error {
	cmd := exec.Command("minikube", "delete")
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AddonsEnable ...
func AddonsEnable(addon string) error {
	cmd := exec.Command("minikube", "addons", "enable", addon)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ConfigSet ...
func ConfigSet(key string, value string) error {
	cmd := exec.Command("minikube", "config", "set", key, value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Version ...
func Version() error {
	cmd := exec.Command("minikube", "version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Ip ...
func Ip() (string, error) {
	out, err := exec.Command("minikube", "ip").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}

// DownloadExecutable ...
func DownloadExecutable(minikubeURL string, minikubeVersion string) error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		fileMap := &download.FileMap{Src: "minikube-windows-amd64.exe", Dst: LOCAL_EXECUTABLE_NAME}
		_, err = download.FromUrl(minikubeURL, minikubeVersion, "minikube", []*download.FileMap{fileMap}, filepath.Dir(localFile))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteExecutable ...
func DeleteExecutable() error {
	localFile := utils.GetBinDir("gokube") + string(os.PathSeparator) + LOCAL_EXECUTABLE_NAME
	return os.RemoveAll(localFile)
}

// DeleteWorkingDirectory ...
func DeleteWorkingDirectory() error {
	return utils.CleanDir(utils.GetUserHome() + string(os.PathSeparator) + ".minikube")
}
