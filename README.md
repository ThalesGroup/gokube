# gokube
[![pages-build-deployment](https://github.com/ThalesGroup/gokube/actions/workflows/pages/pages-build-deployment/badge.svg?branch=master)](https://github.com/ThalesGroup/gokube/actions/workflows/pages/pages-build-deployment) [![Release on github](https://github.com/ThalesGroup/gokube/actions/workflows/github-release.yaml/badge.svg)](https://github.com/ThalesGroup/gokube/actions/workflows/github-release.yaml)

![gokube](https://thalesgroup.github.io/gokube/logo/gokube_150x150.png)

## What is gokube?

gokube is a tool that simplifies day-to-day development with [Kubernetes](https://github.com/kubernetes/kubernetes) on your laptop under Windows.

gokube downloads and installs several dependencies such as:
* [minikube](https://github.com/kubernetes/minikube)
* [docker](https://www.docker.com)
* [helm](https://github.com/helm/helm)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)

You will be able to deploy useful helm charts for development in your Kubernetes cluster with one click.

gokube is configured with a dedicated helm repository named [miniapps](https://thalesgroup.github.io/miniapps) which contains the following charts:
* [cassandra](https://github.com/thalesgroup/miniapps/tree/master/charts/cassandra)
* [mysql](https://github.com/thalesgroup/miniapps/tree/master/charts/mysql)
* [heapster](https://github.com/thalesgroup/miniapps/tree/master/charts/heapster)
* [pact](https://github.com/thalesgroup/miniapps/tree/master/charts/pact)
* [kibana](https://github.com/thalesgroup/miniapps/tree/master/charts/kibana)
* [grafana](https://github.com/thalesgroup/miniapps/tree/master/charts/grafana)

These charts are optimized in terms of memory and CPU for minikube and are very useful for developers.

## How to upgrade gokube?

### Windows

#### Download binary

* The latest release for gokube can be downloaded on the [Releases page](https://github.com/thalesgroup/gokube/releases/latest).
* Copy the executable file to: C:\gokube\bin and replace the previous one.

#### Upgrade gokube

```shell
$ gokube init
```

## How to install gokube?

### Windows

#### Requirements
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads) or [Hyper-V](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md#hyperV-driver)
* VT-x/AMD-v virtualization must be enabled in BIOS
* Internet connection for the first run

#### Assumptions 

You will use C:\gokube\bin to store executable files.

#### Set up your environment

gokube is aware of HTTP_PROXY, HTTPS_PROXY and NO_PROXY environment variables.
When these variables are set, they are used to download the gokube dependencies and to configure the Docker daemon.
You can define different proxy values for the Docker daemon in using the --http-proxy, --https-proxy, and --no-proxy init command flags.

#### Set up your directory

You’ll need a place to store the gokube executable:
* Open Windows Explorer.
* Create a new folder: C:\gokube (assuming you want gokube on your C drive, although this can be placed anywhere).
* Create a subfolder in the gokube folder: C:\gokube\bin.

#### Download binary

* The latest release for gokube can be download on the [Releases page](https://github.com/thalesgroup/gokube/releases/latest).
* Copy executable file to: C:\gokube\bin
* The gokube executable will be named as gokube-<version>-<type>-<platform>.<arch>.exe. Rename the executable to gokube.exe for ease of use.

#### Verify the executable

In your preferred CLI, at the prompt, type gokube and press the Enter key. You should see output that starts with:

```shell
$ gokube
Using environment variable MINIKUBE_MEMORY=12288
Using environment variable MINIKUBE_CPUS=6
gokube is a nice installer to provide an environment for developing day-to-day with kubernetes & helm on your laptop.

Usage:
  gokube [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern + k9s and creates a minikube VM
  pause       Pauses gokube. This command pauses the minikube VM
  reset       Resets gokube. This command restores minikube VM from previously taken snapshot
  resume      Resumes gokube. This command resumes the minikube VM
  save        Creates a gokube reference. This command takes a snapshot of the minikube VM (which will be the reference for reset command)
  start       Starts gokube. This command starts minikube
  stop        Stops gokube. This command stops minikube
  version     Shows version for gokube

Flags:
  -h, --help      help for gokube
  -v, --verbose   Activate verbose logging

Use "gokube [command] --help" for more information about a command.
```
If you do, then the installation is complete.

If you don’t, double-check the path that you placed the gokube.exe file in and ensure you correctly added that path to your PATH variable.

#### Install gokube

```shell
$ gokube init
Using environment variable MINIKUBE_MEMORY=12288
Using environment variable MINIKUBE_CPUS=6
Warning: your Virtualbox GUI shall not be open and no other VM shall be currently running
Press <CTRL+C> within the next 10s if you need to check this or press <ENTER> now to continue...
Deleting previous minikube VM...
Resetting host-only network used by minikube...
Creating minikube VM with kubernetes v1.31.0...
* minikube v1.34.0 on Microsoft Windows 10 Enterprise 10.0.19045.5011 Build 19045.5011
  - MINIKUBE_CPUS=6
  - MINIKUBE_MEMORY=12288
* Using the virtualbox driver based on user configuration
* Starting "minikube" primary control-plane node in "minikube" cluster
* Creating virtualbox VM (CPUs=6, Memory=12288MB, Disk=20480MB) ...
* Found network options:
  - HTTP_PROXY=http://<proxy>:8080
  - HTTPS_PROXY=http://<proxy>:8080
  - NO_PROXY=minikube,127.0.0.1,192.168.99.100
  - http_proxy=http://<proxy>:8080
  - https_proxy=http://<proxy>:8080
  - no_proxy=minikube,127.0.0.1,192.168.99.100
* Preparing Kubernetes v1.17.3 on Docker 19.03.6 ...
  - env HTTP_PROXY=http://<proxy>:8080
  - env HTTPS_PROXY=http://<proxy>:8080
  - env NO_PROXY=minikube,127.0.0.1,192.168.99.100
  - env HTTP_PROXY=http://<proxy>:8080
  - env HTTPS_PROXY=http://<proxy>:8080
  - env NO_PROXY=minikube,127.0.0.1,192.168.99.100
  - Generating certificates and keys ...
  - Booting up control plane ...
  - Configuring RBAC rules ...
* Configuring bridge CNI (Container Networking Interface) ...
  - Using image gcr.io/k8s-minikube/storage-provisioner:v5
* Enabled addons: storage-provisioner, default-storageclass
* Verifying Kubernetes components...
* Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
* dashboard is an addon maintained by Kubernetes. For any concerns contact minikube on GitHub.
You can view the list of minikube maintainers at: https://github.com/kubernetes/minikube/blob/master/OWNERS
  - Using image docker.io/kubernetesui/dashboard:v2.7.0
  - Using image docker.io/kubernetesui/metrics-scraper:v1.0.8
* Some dashboard features require the metrics-server addon. To enable all features please run:

        minikube addons enable metrics-server

* The 'dashboard' addon is enabled
Switched to context "minikube".
Installing ChartMuseum...
"chartmuseum" has been added to your repositories
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "chartmuseum" chart repository
Update Complete. ⎈Happy Helming!⎈
Starting chartmuseum/chartmuseum components...
Waiting for chartmuseum..................
"minikube" has been added to your repositories
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "minikube" chart repository
...Successfully got an update from the "chartmuseum" chart repository
Update Complete. ⎈Happy Helming!⎈
Configuring miniapps repository...
"miniapps" has been added to your repositories
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "minikube" chart repository
...Successfully got an update from the "chartmuseum" chart repository
...Successfully got an update from the "miniapps" chart repository
Update Complete. ⎈Happy Helming!⎈
Exposing kubernetes dashboard to nodeport 30000...
service/kubernetes-dashboard patched

gokube init completed in 3m3s
```

We can see that pods are still being created from the ContainerCreating status:

```shell
$ kubectl get pod --all-namespaces
NAMESPACE              NAME                                        READY   STATUS               RESTARTS       AGE
kube-system            chartmuseum-b9b5d4646-4zq9m                 0/1     ContainerCreating    0              7m36s
kube-system            coredns-6f6b679f8f-8xwrf                    1/1     Running              0              7m38s
kube-system            etcd-minikube                               1/1     Running              0              7m44s
kube-system            kube-apiserver-minikube                     1/1     Running              0              7m44s
kube-system            kube-controller-manager-minikube            1/1     Running              0              7m44s
kube-system            kube-proxy-hmxpz                            1/1     Running              0              7m39s
kube-system            kube-scheduler-minikube                     1/1     Running              0              7m44s
kube-system            storage-provisioner                         1/1     Running              1 (7m8s ago)   7m42s
kubernetes-dashboard   dashboard-metrics-scraper-c5db448b4-rvlxl   0/1     ContainerCreating    0              7m38s
kubernetes-dashboard   kubernetes-dashboard-695b96c756-gf76s       1/1     Running   0              7m38s
```

We can see that pods are now running, and we will now be able to access gokube:

```shell
$ kubectl get pod --all-namespaces
NAMESPACE              NAME                                        READY   STATUS    RESTARTS       AGE
kube-system            chartmuseum-b9b5d4646-4zq9m                 1/1     Running   0              7m36s
kube-system            coredns-6f6b679f8f-8xwrf                    1/1     Running   0              7m38s
kube-system            etcd-minikube                               1/1     Running   0              7m44s
kube-system            kube-apiserver-minikube                     1/1     Running   0              7m44s
kube-system            kube-controller-manager-minikube            1/1     Running   0              7m44s
kube-system            kube-proxy-hmxpz                            1/1     Running   0              7m39s
kube-system            kube-scheduler-minikube                     1/1     Running   0              7m44s
kube-system            storage-provisioner                         1/1     Running   1 (7m8s ago)   7m42s
kubernetes-dashboard   dashboard-metrics-scraper-c5db448b4-rvlxl   1/1     Running   0              7m38s
kubernetes-dashboard   kubernetes-dashboard-695b96c756-gf76s       1/1     Running   0              7m38s
```

We can stop gokube with the following command:
```shell
$ gokube stop
Using environment variable MINIKUBE_MEMORY=12288
Using environment variable MINIKUBE_CPUS=6
Warning: you should not stop a VM with a lot of running pods as the restart will be unstable
Press <CTRL+C> within the next 10s if you need to perform some clean or press <ENTER> now to continue...
Stopping minikube VM...
* Stopping node "minikube"  ...
* 1 node stopped.
```

We can start gokube with the following command:
```shell
$ gokube start
Using environment variable MINIKUBE_MEMORY=12288
Using environment variable MINIKUBE_CPUS=6
Starting minikube VM with kubernetes v1.31.0 and container runtime "docker"...
* minikube v1.34.0 on Microsoft Windows 10 Enterprise 10.0.19045.5011 Build 19045.5011
  - MINIKUBE_CPUS=6
  - MINIKUBE_MEMORY=12288
* Using the virtualbox driver based on existing profile
* Starting "minikube" primary control-plane node in "minikube" cluster
* Restarting existing virtualbox VM for "minikube" ...
* Found network options:
  - HTTP_PROXY=http://<proxy>:8080
  - HTTPS_PROXY=http://<proxy>:8080
  - NO_PROXY=minikube,127.0.0.1,192.168.99.100
  - http_proxy=http://<proxy>:8080
  - https_proxy=http://<proxy>:8080
  - no_proxy=minikube,127.0.0.1,192.168.99.100
* To pull new external images, you may need to configure a proxy: https://minikube.sigs.k8s.io/docs/reference/networking/proxy/
* Preparing Kubernetes v1.31.0 on Docker 27.2.0 ...
  - env HTTP_PROXY=http://<proxy>:8080
  - env HTTPS_PROXY=http://<proxy>:8080
  - env NO_PROXY=minikube,127.0.0.1,192.168.99.100
  - env HTTP_PROXY=http://<proxy>:8080
  - env HTTPS_PROXY=http://<proxy>:8080
  - env NO_PROXY=minikube,127.0.0.1,192.168.99.100
* Configuring bridge CNI (Container Networking Interface) ...
  - Using image gcr.io/k8s-minikube/storage-provisioner:v5
  - Using image docker.io/kubernetesui/dashboard:v2.7.0
  - Using image docker.io/kubernetesui/metrics-scraper:v1.0.8
* Verifying Kubernetes components...
* Some dashboard features require the metrics-server addon. To enable all features please run:

        minikube addons enable metrics-server

* Enabled addons: storage-provisioner, default-storageclass, dashboard
* Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

## Additional links

* [**Contributing**](./CONTRIBUTING.md)
* [**Development Guide**](./docs/developer-guide.md)

