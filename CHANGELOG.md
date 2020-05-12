# GoKube Release Notes

## Version 1.10.0 - 05/12/2020
* Bump to minikube 1.10.0
* Bump to helm-spray 4.0.0 (which implies support only for helm 3)

## Version 1.9.2 - 04/04/2020
* Bump to minikube 1.9.1
* Reduced the timeout to check for new version of gokube

## Version 1.9.1 - 03/31/2020
* Shows a warning if not using the latest gokube release

## Version 1.9.0 - 03/27/2020
* Bump to minikube 1.9.0, K8S 1.18, docker 19.03.8
* Upgrade is now also possible upon restart (when we don't want VM to be respawn)

## Version 1.8.2 - 03/14/2020
* Upgrade is automatically done the first time you execute a new version of gokube
* Miniapps URL repo fixed to match new thalesgroup organization
* Bump to kubernetes v1.17.3, helm 2.16.3, minikube v1.8.2

## Version 1.8.1 - 02/08/2020
* Bump to kubernetes v1.17.2, minikube v1.7.2

## Version 1.8.0 - 12/19/2019
THE BIG ONE

* Bump to kubernetes v1.16.4, minikube v1.6.1, helm v2.16.1 and docker 19.03.3
* Support of Virtualbox 6 (with the limitation that no other VMs shall be running during gokube init)
* Added warning messages on gokube stop to prevent crashes/unstabilities on VM restart
* Support of environment variables to download more recent versions of dependencies (to avoid generating a new version of gokube each time)
* Added pause and resume commands

Supported environment variables are:
- KUBERNETES_VERSION (ex: "v1.16.4")
- MINIKUBE_VERSION (ex: "v1.6.1")
- DOCKER_VERSION (ex: "19.03.3")
- HELM_VERSION (ex: "v2.16.1")
- HELM_SPRAY_VERSION (ex: "v3.4.5")

As recent minikube versions automatically selects the best hypervisor, please note no tests have been made yet for Win10/Hyper-V contexts

## Version 1.7.7 - 10/16/2019
* Update helm-spray to 3.4.5 [#42](https://github.com/gemalto/helm-spray/pull/42)

## Version 1.7.6 - 09/23/2019
* Fix for 1.7.5 "Unable to read config file" [#13](https://github.com/gemalto/gokube/issues/13)

## Version 1.7.5 - 09/18/2019
* Enhance management of default kubernetes version (hotfix)

## Version 1.7.4 - 09/12/2019
* Enhance management of default kubernetes version
<p>Changing the kubernetes version with the environment variable did not prevent a potential upgrade of kubernetes
upon a VM restart (e.g. gokube stop / gokube start). The desired version of kubernetes is now stored under
.gokube/config.yaml and is used for init and start commands</p>

## Version 1.7.3 - 09/11/2019
* Enhance management of default kubernetes version
<p>You can change the default kubernetes version (v1.10.13) for your VM in setting a KUBERNETES_VERSION global environment variable</p>

## Version 1.7.2 - 09/10/2019
* Bump to helm spray v3.4.4 (which fixes issues on liveness/readiness for id-provider) [#12](https://github.com/gemalto/gokube/pull/12)

## Version 1.7.1 - 08/14/2019
* Bump to minikube v1.3.1 (which fixes a TTY/PTY issue preventing gokube init progress bar to be displayed)

## Version 1.7.0 - 08/06/2019
* Bump to minikube v1.3.0 (to benefit from the NAT DNS options), bump to helm spray v3.4.3 [#10](https://github.com/gemalto/gokube/pull/10)

## Version 1.6.2 - 05/29/2019
* kubernetes-dashboard not always patched [#9](https://github.com/gemalto/gokube/pull/9)

## Version 1.6.1 - 05/22/2019
* Fix for kubernetes-dashboard not always patched [#8](https://github.com/gemalto/gokube/issues/8)

## Version 1.6.0 - 04/02/2019
* Bump to minikube v1.0.0 [#7](https://github.com/gemalto/gokube/pull/7)

## Version 1.5.0 - 02/08/2019
* Bump to helm-spray v3.2.0 and stern v1.10.0 [#6](https://github.com/gemalto/gokube/pull/6)

## Version 1.4.0 - 11/21/2018
* Bump to monocular v1.2.0 (that support http proxy)
* Bump to docker v18.06
* Bump to k8s v1.10.9
* Add docker image 'mongodb' in cache for monocular
* Add support of any-proxy (for transparent proxy)
* Add support of --cache-images
* Patch kubernetes-dashboard service to expose it on nodePort: 30000

## Version 1.3.0 - 10/15/2018
* Add option for using a specific version for Minikube
* Add option for using a specific version for Helm
* Add option for using a specific fork for Minikube

## Version 1.2.0 - 10/09/2018
* Updated Minikube to 0.30.0
* Updated Kubectl to 1.12.0
* Remove directories: '.kube' and '.docker' before upgrading gokube
