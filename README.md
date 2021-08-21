# gokube
[![Build Status](https://api.travis-ci.org/thalesgroup/gokube.svg?branch=master)](https://travis-ci.org/thalesgroup/gokube)

![gokube](https://thalesgroup.github.io/gokube/logo/gokube_150x150.png)

## What is gokube?

gokube is a tool that makes it easy developing day-to-day with [Kubernetes](https://github.com/kubernetes/kubernetes) on your laptop under Windows.

gokube downloads and installs many dependencies such as:
* [minikube](https://github.com/kubernetes/minikube)
* [docker](https://www.docker.com)
* [helm](https://github.com/helm/helm)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [monocular](https://github.com/helm/monocular)

gokube deploys and configures Monocular for a better user experience!
You will be able to deploy in one click useful helm charts for developing in your kubernetes cluster.

gokube is configured with a dedicated helm repository named [miniapps](https://thalesgroup.github.io/miniapps) which contains the following charts:
* [cassandra](https://github.com/thalesgroup/miniapps/tree/master/charts/cassandra)
* [mysql](https://github.com/thalesgroup/miniapps/tree/master/charts/mysql)
* [heapster](https://github.com/thalesgroup/miniapps/tree/master/charts/heapster)
* [pact](https://github.com/thalesgroup/miniapps/tree/master/charts/pact)
* [kibana](https://github.com/thalesgroup/miniapps/tree/master/charts/kibana)
* [grafana](https://github.com/thalesgroup/miniapps/tree/master/charts/grafana)

These charts are optimized in term of memory and cpu for minikube and very useful for developers.

## Continuous Integration & Delivery 
[![Build Status](https://api.travis-ci.org/thalesgroup/gokube.svg?branch=master)](https://travis-ci.org/thalesgroup/gokube)

gokube is building and delivering under Travis and Github Actions.

## How to upgrade gokube?

### Windows

#### Download binary

* The latest release for gokube can be download on the [Releases page](https://github.com/thalesgroup/gokube/releases/latest).
* Copy executable file to: C:\gokube\bin and replace the previous one.

#### Upgrade gokube

```shell
$ gokube init
```

## How to install gokube?

### Windows

#### Requirements
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads) or [Hyper-V](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md#hyperV-driver)
* VT-x/AMD-v virtualization must be enabled in BIOS
* Internet connection on first run

#### Assumptions 

You will use C:\gokube\bin to store executable files.

#### Set up your environment

gokube is aware of HTTP_PROXY, HTTPS_PROXY and NO_PROXY environment variables.
When these variables are set, they are used to download the gokube dependencies and to configure docker daemon.
You can define different proxy values for docker daemon in using --http-proxy, --https-proxy and --no-proxy init command flags

#### Set up your directory

You’ll need a place to store the gokube executable:
* Open Windows Explorer.
* Create a new folder: C:\gokube, assuming you want gokube on your C drive, although this can go anywhere.
* Create a subfolder in the gokube folder: C:\gokube\bin

#### Download binary

* The latest release for gokube can be download on the [Releases page](https://github.com/thalesgroup/gokube/releases/latest).
* Copy executable file to: C:\gokube\bin
* The gokube executable will be named as gokube-version-type+platform.arch.exe. Rename the executable to gokube.exe for ease of use.

#### Verify the executable

In your preferred CLI, at the prompt, type gokube and press the Enter key. You should see output that starts with:

```shell
$ gokube
gokube is a nice installer to provide an environment for developing day-to-day with kubernetes & helm on your laptop.

Usage:
  gokube [command]

Available Commands:
  help        Help about any command
  init        Initializes gokube. This command downloads dependencies: minikube + helm + kubectl + docker + stern and creates the virtual machine (minikube)
  pause       Pauses minikube. This command pauses minikube VM
  resume      Resumes minikube. This command resumes minikube VM
  start       Starts minikube. This command starts minikube
  stop        Stops minikube. This command stops minikube
  version     Shows version for gokube

Flags:
  -h, --help   help for gokube

Use "gokube [command] --help" for more information about a command.
```
If you do, then the installation is complete.

If you don’t, double-check the path that you placed the gokube.exe file in and that you typed that path correctly when you added it to your PATH variable.

#### Install gokube

```shell
$ gokube init
WARNING: your Virtualbox GUI shall not be open and no other VM shall be currently running
Press <CTRL+C> within the next 10s it you need to check this or press <ENTER> now to continue...
Deleting previous minikube VM...
Resetting host-only network used by minikube...
This version of gokube is launched for the first time, forcing upgrade...
Downloading gokube dependencies...
minikube v1.8.2: 52.12 MiB / 52.12 MiB [--------------------------------------] 100.00% 1.63 MiB p/s
helm v2.16.3: 24.16 MiB / 24.16 MiB [-----------------------------------------] 100.00% 1.81 MiB p/s
docker v19.03.3: 65.28 MiB / 65.28 MiB [--------------------------------------] 100.00% 1.39 MiB p/s
kubectl v1.17.3: 41.95 MiB / 41.95 MiB [--------------------------------------] 100.00% 2.23 MiB p/s
stern v1.11.0: 20.93 MiB / 20.93 MiB [----------------------------------------] 100.00% 1.36 MiB p/s
Creating minikube VM with kubernetes v1.17.3...
* minikube v1.8.2 on Microsoft Windows 7 Enterprise Service Pack 1 6.1.7601 Build 7601
* Automatically selected the virtualbox driver
* Downloading VM boot image ...
    > minikube-v1.8.0.iso.sha256: 65 B / 65 B [--------------] 100.00% ? p/s 0s
    > minikube-v1.8.0.iso: 173.56 MiB / 173.56 MiB [ 100.00% 2.23 MiB p/s 1m18s
* Downloading preloaded images tarball for k8s v1.17.3 ...
    > preloaded-images-k8s-v1-v1.17.3-docker-overlay2.tar.lz4: 499.26 MiB / 499
* Creating virtualbox VM (CPUs=4, Memory=8192MB, Disk=20000MB) ...
* Found network options:
  - HTTP_PROXY=http://10.43.216.8:8080
  - HTTPS_PROXY=http://10.43.216.8:8080
  - NO_PROXY=minikube,dockerhub.gemalto.com,r-buildggs.gemalto.com,nexusfreel1-emea-proxy.gemalto.com,127.0.0.1,192.168.99.100
  - http_proxy=http://10.43.216.8:8080
  - https_proxy=http://10.43.216.8:8080
  - no_proxy=minikube,dockerhub.gemalto.com,r-buildggs.gemalto.com,nexusfreel1-emea-proxy.gemalto.com,127.0.0.1,192.168.99.100
* Preparing Kubernetes v1.17.3 on Docker 19.03.6 ...
  - env HTTP_PROXY=http://10.43.216.8:8080
  - env HTTPS_PROXY=http://10.43.216.8:8080
  - env NO_PROXY=minikube,dockerhub.gemalto.com,r-buildggs.gemalto.com,nexusfreel1-emea-proxy.gemalto.com,127.0.0.1,192.168.99.100
  - env HTTP_PROXY=http://10.43.216.8:8080
  - env HTTPS_PROXY=http://10.43.216.8:8080
  - env NO_PROXY=minikube,dockerhub.gemalto.com,r-buildggs.gemalto.com,nexusfreel1-emea-proxy.gemalto.com,127.0.0.1,192.168.99.100
  - apiserver.runtime-config=apps/v1beta1=true,apps/v1beta2=true,extensions/v1beta1/daemonsets=true,extensions/v1beta1/deployments=true,extensions/v1beta1/replicasets=true,extensions/v1beta1/networkpolicies=true,extensions/v1beta1/podsecuritypolicies=true
* Launching Kubernetes ...
* Enabling addons: default-storageclass, storage-provisioner
* Waiting for cluster to come online ...
* Done! kubectl is now configured to use "minikube"
* The 'dashboard' addon is enabled
Caching additional docker images...
Switched to context "minikube".
Initializing helm...
Creating C:\Users\user\.helm
Creating C:\Users\user\.helm\repository
Creating C:\Users\user\.helm\repository\cache
Creating C:\Users\user\.helm\repository\local
Creating C:\Users\user\.helm\plugins
Creating C:\Users\user\.helm\starters
Creating C:\Users\user\.helm\cache\archive
Creating C:\Users\user\.helm\repository\repositories.yaml
Adding stable repo with URL: https://kubernetes-charts.storage.googleapis.com
Adding local repo with URL: http://127.0.0.1:8879/charts
$HELM_HOME has been configured at C:\Users\user\.helm.

Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.

Please note: by default, Tiller is deployed with an insecure 'allow unauthenticated users' policy.
To prevent this, run `helm init` with the --tiller-tls-verify flag.
For more information on securing your installation see: https://docs.helm.sh/using_helm/#securing-your-helm-installation
"miniapps" has been added to your repositories
Hang tight while we grab the latest from your chart repositories...
...Skip local chart repository
...Successfully got an update from the "miniapps" chart repository
...Successfully got an update from the "stable" chart repository
Update Complete.
helm-spray v3.4.5: 5.02 MiB / 5.02 MiB [------------------------------------] 100.00% 631.12 KiB p/s
Exposing kubernetes dashboard to nodeport 30000...
service/kubernetes-dashboard patched

gokube has been installed.
```

We can see that pods are still being created from the ContainerCreating status:

```shell
$ kubectl get pod --all-namespaces
NAMESPACE    NAME                                                 READY  STATUS             RESTARTS AGE
kube-system  etcd-minikube                                        1/1    Running            0        1m
kube-system  gokube-mongodb-7c86445c7-zx6cv                       0/1    ContainerCreating  0        1m
kube-system  gokube-monocular-api-5798749fdb-6xwsn                0/1    ContainerCreating  0        1m
kube-system  gokube-monocular-prerender-c9f57f6c8-nbjl7           0/1    ContainerCreating  0        1m
kube-system  gokube-monocular-ui-7d79f486-w98px                   0/1    ContainerCreating  0        1m
kube-system  kube-addon-manager-minikube                          1/1    Running            0        1m
kube-system  kube-apiserver-minikube                              1/1    Running            0        1m
kube-system  kube-controller-manager-minikube                     1/1    Running            0        1m
kube-system  kube-dns-86f4d74b45-4swsw                            3/3    Running            0        2m
kube-system  kube-proxy-dltxx                                     1/1    Running            0        2m
kube-system  kube-scheduler-minikube                              1/1    Running            0        1m
kube-system  kubernetes-dashboard-5498ccf677-5p82h                1/1    Running            0        2m
kube-system  nginx-nginx-ingress-controller-859558948c-5rr2c      1/1    Running            0        1m
kube-system  nginx-nginx-ingress-default-backend-7bb66746b9-tfmm2 1/1    Running            0        1m
kube-system  storage-provisioner                                  1/1    Running            0        2m
kube-system  tiller-deploy-f9b8476d-rk5ps                         1/1    Running            0        2m
```

We can see that pods are now running and we will now be able to access to gokube:

```shell
$ kubectl get pod --all-namespaces
NAMESPACE    NAME                                                 READY  STATUS   RESTARTS AGE
kube-system  etcd-minikube                                        1/1    Running  0        5m
kube-system  gokube-mongodb-7c86445c7-zx6cv                       1/1    Running  0        6m
kube-system  gokube-monocular-api-5798749fdb-6xwsn                1/1    Running  4        6m
kube-system  gokube-monocular-prerender-c9f57f6c8-nbjl7           1/1    Running  0        6m
kube-system  gokube-monocular-ui-7d79f486-w98px                   1/1    Running  0        6m
kube-system  kube-addon-manager-minikube                          1/1    Running  0        5m
kube-system  kube-apiserver-minikube                              1/1    Running  0        5m
kube-system  kube-controller-manager-minikube                     1/1    Running  0        5m
kube-system  kube-dns-86f4d74b45-4swsw                            3/3    Running  0        6m
kube-system  kube-proxy-dltxx                                     1/1    Running  0        6m
kube-system  kube-scheduler-minikube                              1/1    Running  0        5m
kube-system  kubernetes-dashboard-5498ccf677-5p82h                1/1    Running  0        6m
kube-system  nginx-nginx-ingress-controller-859558948c-5rr2c      1/1    Running  0        6m
kube-system  nginx-nginx-ingress-default-backend-7bb66746b9-tfmm2 1/1    Running  0        6m
kube-system  storage-provisioner                                  1/1    Running  0        6m
kube-system  tiller-deploy-f9b8476d-rk5ps                         1/1    Running  0        6m
```

We can stop gokube running the following command:
```shell
$ gokube stop
Stopping local Kubernetes cluster...
Machine stopped.
```

We can start gokube running the following command:
```shell
$ gokube start
Starting local Kubernetes v1.10.0 cluster...
Starting VM...
Getting VM IP address...
Moving files into cluster...
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
Loading cached images from config file.
```

## Additional links

* [**Contributing**](./CONTRIBUTING.md)
* [**Development Guide**](./docs/developer-guide.md)

