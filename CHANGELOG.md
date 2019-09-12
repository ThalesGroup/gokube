# GoKube Release Notes

## Version 1.7.4 - 09/12/2019
* Enhance management of default kubernetes version (hotfix)
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
