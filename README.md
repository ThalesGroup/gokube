# GoKube

[![Build Status](https://api.travis-ci.com/gemalto/gokube.svg?branch=master)](https://travis-ci.com/gemalto/gokube)

## What is GoKube?

![gokube](https://github.com/gemalto/gokube/blob/master/logo/gokube_150x150.png)

GoKube is a tool that makes it easy developping day-to-day with [Kubernetes](https://github.com/kubernetes/kubernetes) on your laptop under Windows.

GoKube downloads and installs many dependencies such as:
* [Minikube](https://github.com/kubernetes/minikube)
* [Docker](https://www.docker.com)
* [Helm](https://github.com/helm/helm)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [Monocular](https://github.com/helm/monocular)

GoKube deploys and configures Monocular for a better user experience!
You will be able to deploy in one click useful helm charts for developing in your kubernetes cluster.

GoKube is configured with a dedicated helm repository named [miniapps](https://github.com/gokube/miniapps) which contains the following charts:
* [Cassandra](https://github.com/gemalto/miniapps/tree/master/charts/cassandra)
* [Heapster](https://github.com/gemalto/miniapps/tree/master/charts/heapster)
* [Pact](https://github.com/gemalto/miniapps/tree/master/charts/pact)
* [Kibana](https://github.com/gemalto/miniapps/tree/master/charts/kibana)
* [Grafana](https://github.com/gemalto/miniapps/tree/master/charts/grafana)
* [Livedoc](https://github.com/gemalto/miniapps/tree/master/charts/livedoc)

These charts are optimized in term of memory and cpu for minikube and very useful for developers.

## Requirements
* Windows
    * [VirtualBox](https://www.virtualbox.org/wiki/Downloads) or [Hyper-V](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md#hyperV-driver)
* VT-x/AMD-v virtualization must be enabled in BIOS
* Internet connection on first run

## How to install GoKube?

### Windows


#### Assumptions 

You will use C:\gokube\bin to store executable files.

#### Set up Your Directory

You’ll need a place to store the gokube executable:
* Open Windows Explorer.
* Create a new folder: C:\gokube, assuming you want gokube on your C drive, although this can go anywhere.
* Create a subfolder in the gokube folder: C:\gokube\bin

#### Download binary

* The latest release for gokube can be download on the [Releases page](https://github.com/gemalto/gokube/releases/latest).
* Copy executable file to: C:\gokube\bin
* The gokube executable will be named as gokube-version-type+platform.arch.exe. Rename the executable to gokube.exe for ease of use.

#### Verify the Executable

In your preferred CLI, at the prompt, type gokube and press the Enter key. You should see output that starts with:

```shell
$ gokube

gokube is a nice installer to provide an environment for developing day-to-day 
with kubernetes & helm on your laptop.

Usage:
  gokube [command]

Available Commands:
  help        Help about any command
  init        Initializes gokube. This command downloads dependencies: 
              minikube + helm + kubectl + docker + monocular and creates the virtual machine (minikube)
  version     Shows version for gokube

Flags:
  -h, --help   help for gokube
```
If you do, then the installation is complete.

If you don’t, double-check the path that you placed the gokube.exe file in and that you typed that path correctly when you added it to your PATH variable.

## Quickstart
Here's a brief demo of GoKube usage.

```shell
$ gokube init

minikube v0.28.0: 40.83 MiB / 40.83 MiB [-------------------------------------] 100.00% 2.21 MiB p/s
helm v2.9.1: 8.78 MiB / 8.78 MiB [--------------------------------------------] 100.00% 2.20 MiB p/s
docker v17.09.0: 16.17 MiB / 16.17 MiB [--------------------------------------] 100.00% 2.11 MiB p/s
kubectl v1.10.0: 52.16 MiB / 52.16 MiB [--------------------------------------] 100.00% 1.61 MiB p/s

Installing goKube!

Starting local Kubernetes v1.10.0 cluster...
Starting VM...
Downloading Minikube ISO
 153.08 MB / 153.08 MB  100.00% 0ssss
Getting VM IP address...
Waiting for image caching to complete...
Moving files into cluster...
Downloading kubelet v1.10.0
Downloading kubeadm v1.10.0
Finished Downloading kubeadm v1.10.0
Finished Downloading kubelet v1.10.0
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
Loading cached images from config file.
Switched to context "minikube".
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
Happy Helming!
"monocular" has been added to your repositories
Hang tight while we grab the latest from your chart repositories...
...Skip local chart repository
...Successfully got an update from the "monocular" chart repository
...Successfully got an update from the "stable" chart repository
Update Complete. ⎈ Happy Helming!⎈
Starting stable/nginx-ingress components...
Starting monocular/monocular components...

goKube! has been installed.

To verify that goKube! has started, run:
> kubectl get pods --all-namespaces
```

We can see that pods are still being created from the ContainerCreating status:

```shell
$ kubectl get pod --all-namespaces
NAMESPACE     NAME                                                   READY     STATUS              RESTARTS   AGE
kube-system   etcd-minikube                                          1/1       Running             0          1m
kube-system   gokube-mongodb-7c86445c7-zx6cv                         0/1       ContainerCreating   0          1m
kube-system   gokube-monocular-api-5798749fdb-6xwsn                  0/1       ContainerCreating   0          1m
kube-system   gokube-monocular-prerender-c9f57f6c8-nbjl7             0/1       ContainerCreating   0          1m
kube-system   gokube-monocular-ui-7d79f486-w98px                     0/1       ContainerCreating   0          1m
kube-system   kube-addon-manager-minikube                            1/1       Running             0          1m
kube-system   kube-apiserver-minikube                                1/1       Running             0          1m
kube-system   kube-controller-manager-minikube                       1/1       Running             0          1m
kube-system   kube-dns-86f4d74b45-4swsw                              3/3       Running             0          2m
kube-system   kube-proxy-dltxx                                       1/1       Running             0          2m
kube-system   kube-scheduler-minikube                                1/1       Running             0          1m
kube-system   kubernetes-dashboard-5498ccf677-5p82h                  1/1       Running             0          2m
kube-system   nginx-nginx-ingress-controller-859558948c-5rr2c        1/1       Running             0          1m
kube-system   nginx-nginx-ingress-default-backend-7bb66746b9-tfmm2   1/1       Running             0          1m
kube-system   storage-provisioner                                    1/1       Running             0          2m
kube-system   tiller-deploy-f9b8476d-rk5ps                           1/1       Running             0          2m
```

We can see that pods are now running and we will now be able to access to gokube:

```shell
$ kubectl get pod --all-namespaces
NAMESPACE     NAME                                                   READY     STATUS    RESTARTS   AGE
kube-system   etcd-minikube                                          1/1       Running   0          5m
kube-system   gokube-mongodb-7c86445c7-zx6cv                         1/1       Running   0          6m
kube-system   gokube-monocular-api-5798749fdb-6xwsn                  1/1       Running   4          6m
kube-system   gokube-monocular-prerender-c9f57f6c8-nbjl7             1/1       Running   0          6m
kube-system   gokube-monocular-ui-7d79f486-w98px                     1/1       Running   0          6m
kube-system   kube-addon-manager-minikube                            1/1       Running   0          5m
kube-system   kube-apiserver-minikube                                1/1       Running   0          5m
kube-system   kube-controller-manager-minikube                       1/1       Running   0          5m
kube-system   kube-dns-86f4d74b45-4swsw                              3/3       Running   0          6m
kube-system   kube-proxy-dltxx                                       1/1       Running   0          6m
kube-system   kube-scheduler-minikube                                1/1       Running   0          5m
kube-system   kubernetes-dashboard-5498ccf677-5p82h                  1/1       Running   0          6m
kube-system   nginx-nginx-ingress-controller-859558948c-5rr2c        1/1       Running   0          6m
kube-system   nginx-nginx-ingress-default-backend-7bb66746b9-tfmm2   1/1       Running   0          6m
kube-system   storage-provisioner                                    1/1       Running   0          6m
kube-system   tiller-deploy-f9b8476d-rk5ps                           1/1       Running   0          6m
```

We can stop gokube running the following command:
```shell
$ minikube stop
Stopping local Kubernetes cluster...
Machine stopped.
```

We can start gokube running the following command:
```shell
$ minikube start
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

## Developer Guide

If you want to contribute to this project you are encouraged to send issue request, or provide pull-requests. 
Please read the [developer guide](./docs/developer-guide.md) to learn more on how you can contribute. 
